package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/common"
	"github.com/laghoule/kratos/pkg/config"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
)

// Deployment contain the kubernetes clientset and configuration of the release
type Deployment struct {
	Clientset kubernetes.Interface
	*config.Config
}

// checkOwnership check if it's safe to create, update or delete the deployment
func (d *Deployment) checkOwnership(name, namespace string) error {
	dep, err := d.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting deployment failed: %s", err)
	}

	// managed by kratos
	if err := checkKratosManaged(dep.Labels); err == nil {
		if dep.Labels[DepLabelName] == name {
			return nil
		}
	}

	return fmt.Errorf("deployment is not managed by kratos")
}

// CreateUpdate create or update a deployment
func (d *Deployment) CreateUpdate(name, namespace string) error {
	if err := d.checkOwnership(name, namespace); err != nil {
		return err
	}

	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return fmt.Errorf("converting label failed: %s", err)
	}

	// merge common & deployments labels
	if err := mergo.Map(&d.Deployment.Labels, d.Common.Labels); err != nil {
		return fmt.Errorf("merging deployment labels failed: %s", err)
	}

	// deep copy of map for podLabels
	podLabels := map[string]string{}
	if err := copier.Copy(&podLabels, &d.Deployment.Labels); err != nil {
		return fmt.Errorf("copying deployment labels values failed: %s", err)
	}

	// merge kratosLabels & deployment labels
	if err := mergo.Map(&d.Deployment.Labels, map[string]string(kratosLabel)); err != nil {
		return fmt.Errorf("merging deployment labels failed: %s", err)
	}

	// merge common & deployments annotations
	if err := mergo.Map(&d.Deployment.Annotations, d.Common.Annotations); err != nil {
		return fmt.Errorf("merging deployment annotations failed: %s", err)
	}

	containers := []corev1.Container{}
	volumesMount := []corev1.VolumeMount{}
	volumes := []corev1.Volume{}

	for _, container := range d.Deployment.Containers {
		volumesMount, volumes = getVolumesConfForContainer(name, &container, d.Config)

		containers = append(containers, corev1.Container{
			Name:  container.Name,
			Image: container.Image + ":" + container.Tag,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: d.Deployment.Port,
				},
			},
			Resources:      container.FormatResources(),
			LivenessProbe:  container.FormatProbe(config.LiveProbe),
			ReadinessProbe: container.FormatProbe(config.ReadyProbe),
			VolumeMounts:   volumesMount,
		})
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: d.Deployment.Annotations,
			Labels: labels.Merge(
				d.Deployment.Labels,
				labels.Set{
					DepLabelName: name,
				},
			),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &d.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					DepLabelName: name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        name,
					Namespace:   namespace,
					Annotations: d.Deployment.Annotations,
					Labels: labels.Merge(
						podLabels,
						labels.Set{
							DepLabelName: name,
						},
					),
				},
				Spec: corev1.PodSpec{
					Containers:                   containers,
					AutomountServiceAccountToken: common.BoolPTR(false),
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: common.BoolPTR(true),
					},
					Volumes: volumes,
				},
			},
		},
	}

	_, err = d.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err = d.Clientset.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("updating deployment failed: %s", err)
			}
		} else {
			return fmt.Errorf("creating deployment failed: %s", err)
		}
	}

	return nil
}

// Delete the specified deployment
func (d *Deployment) Delete(name, namespace string) error {
	if err := d.checkOwnership(name, namespace); err != nil {
		return err
	}

	if err := d.Clientset.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting deployment failed: %s", err)
	}

	return nil
}

// List the deployments of the specified namespace
func (d *Deployment) List(namespace string) ([]appsv1.Deployment, error) {
	list, err := d.Clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: config.ManagedLabel,
	})
	if err != nil {
		return nil, fmt.Errorf("getting deployments list failed: %s", err)
	}

	return list.Items, nil
}
