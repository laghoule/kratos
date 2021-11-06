package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/common"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateUpdateService create or update a service
func (c *Client) CreateUpdateService(name, namespace string, port int32) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(common.DeployLabel)
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
					appLabelName: name,
				},
			),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: name,
					Port: port,
					TargetPort: intstr.IntOrString{
						IntVal: port,
					},
				},
			},
			Selector: map[string]string{
				appLabelName: name,
			},
		},
	}

	_, err = c.Clientset.CoreV1().Services(namespace).Create(context.Background(), svc, metav1.CreateOptions{})
	if err != nil {
		// if service exist, we call update
		if errors.IsAlreadyExists(err) {
			if err := c.updateService(name, namespace, port); err != nil {
				return fmt.Errorf("updating service failed: %s", err)
			}
		} else {
			return fmt.Errorf("creating service failed: %s", err)
		}
	}

	return nil
}

// updateService update an existing service
func (c *Client) updateService(name, namespace string, port int32) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(common.DeployLabel)
	if err != nil {
		return nil
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
				kratosLabel,
				labels.Set{
					appLabelName: name,
				},
			),
			ResourceVersion: svcInfo.ResourceVersion,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: svcInfo.Spec.ClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: name,
					Port: port,
					TargetPort: intstr.IntOrString{
						IntVal: port,
					},
				},
			},
			Selector: map[string]string{
				appLabelName: name,
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
	if err := c.Clientset.CoreV1().Services(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("delete service failed: %s", err)
	}

	return nil
}
