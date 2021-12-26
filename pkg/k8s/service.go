package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
)

// checkServiceOwnership check if it's safe to create, update or delete the service
func (c *Client) checkServiceOwnership(name, namespace string) error {
	svc, err := c.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting service failed: %s", err)
	}

	// owned by the kratos release
	if svc.Labels[DepLabelName] == name {
		return nil
	}

	// managed by kratos
	if err := checkKratosManaged(svc.Labels); err == nil {
		return nil
	}

	return fmt.Errorf("service is not managed by kratos")
}

// CreateUpdateService create or update a service
func (c *Client) CreateUpdateService(name, namespace string, conf *config.Config) error {
	if err := c.checkServiceOwnership(name, namespace); err != nil {
		return err
	}

	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					DepLabelName: name,
				},
			),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: name,
					Port: conf.Deployment.Port,
					TargetPort: intstr.IntOrString{
						IntVal: conf.Deployment.Port,
					},
				},
			},
			Selector: map[string]string{
				DepLabelName: name,
			},
		},
	}

	_, err = c.Clientset.CoreV1().Services(namespace).Create(context.Background(), svc, metav1.CreateOptions{})
	if err != nil {
		// if service exist, we call update
		if errors.IsAlreadyExists(err) {
			if err := c.updateService(name, namespace, conf); err != nil {
				return fmt.Errorf("updating service failed: %s", err)
			}
		} else {
			return fmt.Errorf("creating service failed: %s", err)
		}
	}

	return nil
}

// updateService update an existing service. Used by CreateUpdateService.
func (c *Client) updateService(name, namespace string, conf *config.Config) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	// deep copy of map for svcLabels
	svcLabels := map[string]string{}
	if err := copier.Copy(&svcLabels, &conf.Common.Labels); err != nil {
		return fmt.Errorf("copying common labels values failed: %s", err)
	}

	// TODO create a func for merge

	// merge kratosLabel & service labels
	if err := mergo.Map(&svcLabels, kratosLabel); err != nil {
		return fmt.Errorf("merging common labels failed: %s", err)
	}

	// merge common & service labels
	if err := mergo.Map(&svcLabels, conf.Common.Labels); err != nil {
		return fmt.Errorf("merging common labels failed: %s", err)
	}

	svcInfo, err := c.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("getting service informations failed: %s", err)
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				svcLabels,
				labels.Set{
					DepLabelName: name,
				},
			),
			Annotations:     conf.Common.Annotations,
			ResourceVersion: svcInfo.ResourceVersion,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: svcInfo.Spec.ClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: name,
					Port: conf.Deployment.Port,
					TargetPort: intstr.IntOrString{
						IntVal: conf.Deployment.Port,
					},
				},
			},
			Selector: map[string]string{
				DepLabelName: name,
			},
		},
	}

	_, err = c.Clientset.CoreV1().Services(namespace).Update(context.Background(), svc, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("updating service failed: %s", err)
	}

	return nil
}

// DeleteService delete the specified service
func (c *Client) DeleteService(name, namespace string) error {
	if err := c.checkServiceOwnership(name, namespace); err != nil {
		return err
	}

	if err := c.Clientset.CoreV1().Services(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("delete service failed: %s", err)
	}

	return nil
}
