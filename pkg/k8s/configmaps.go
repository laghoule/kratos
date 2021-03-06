package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// ConfigMaps is the interface for configMaps
type ConfigMaps interface {
	CreateUpdate(string, string) error
	Delete(string, string) error
	List(string) ([]corev1.ConfigMap, error)
}

// configMaps contain the kubernetes clientset and configuration of the release
type configMaps struct {
	Clientset kubernetes.Interface
	*config.Config
}

// checkOwnership check if it's safe to create, update or delete the configmaps
func (c *configMaps) checkOwnership(name, namespace string) error {
	configmap, err := c.Clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting configmap failed: %s", err)
	}

	// managed by kratos
	if err := checkKratosManaged(configmap.Labels); err == nil {
		if configmap.Labels[ConfigMapsLabelName] == name {
			return nil
		}
	}

	return fmt.Errorf("configmap is not managed by kratos")
}

// CreateUpdate create or update a configmaps
func (c *configMaps) CreateUpdate(name, namespace string) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return fmt.Errorf("converting label failed: %s", err)
	}

	// merge labels
	if c.Common != nil && c.Common.Labels != nil {
		if err := mergeStringMaps(&c.ConfigMaps.Labels, c.Common.Labels, kratosLabel); err != nil {
			return fmt.Errorf("merging configmaps labels failed: %s", err)
		}
	}

	// merge annotations
	if c.Common != nil && c.Common.Annotations != nil {
		if err := mergeStringMaps(&c.ConfigMaps.Annotations, c.Common.Annotations); err != nil {
			return fmt.Errorf("merging configmaps annotations failed: %s", err)
		}
	}

	for _, file := range c.ConfigMaps.Files {
		cmName := name + "-" + file.Name

		configmaps := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: namespace,
				Labels: labels.Merge(
					c.ConfigMaps.Labels,
					labels.Set{
						ConfigMapsLabelName: cmName,
					},
				),
				Annotations: c.Common.Annotations,
			},
			Data: map[string]string{
				file.Name: file.Data,
			},
		}

		if err := c.checkOwnership(cmName, namespace); err != nil {
			return err
		}

		_, err = c.Clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), configmaps, metav1.CreateOptions{})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				_, err = c.Clientset.CoreV1().ConfigMaps(namespace).Update(context.Background(), configmaps, metav1.UpdateOptions{})
				if err != nil {
					return fmt.Errorf("updating configmaps failed: %s", err)
				}
			} else {
				return fmt.Errorf("creating configmaps failes: %s", err)
			}
		}
	}

	return nil
}

// Delete the configmaps contained in conf for the specified namespace
func (c *configMaps) Delete(name, namespace string) error {
	for _, file := range c.ConfigMaps.Files {
		if err := c.delete(name+"-"+file.Name, namespace); err != nil {
			return err
		}
	}

	return nil
}

// delete a configmaps from a namespace
func (c *configMaps) delete(name, namespace string) error {
	if err := c.checkOwnership(name, namespace); err != nil {
		return err
	}

	if err := c.Clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting configmaps %s failed: %s", name, err)
	}

	return nil
}

// Get a configmap from a namespace
func (c *configMaps) get(name, namespace string) (*corev1.ConfigMap, error) {
	configmap, err := c.Clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting configmap %s failed: %s", name, err)
	}

	return configmap, nil
}

// List the configmaps in the specified namespace
func (c *configMaps) List(namespace string) ([]corev1.ConfigMap, error) {
	list, err := c.Clientset.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: config.ManagedLabel,
	})
	if err != nil {
		return nil, fmt.Errorf("getting configmaps list failed: %s", err)
	}

	return list.Items, nil
}
