package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// ListDeployments list deployment of k8dep labels
func (c *Client) ListDeployments(namespace string) ([]appsv1.Deployment, error) {
	list, err := c.Clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: config.DeployLabel,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting list of deployment: %s", err)
	}

	return list.Items, nil
}

// CreateUpdateDeployment create or update a deployment
func (c *Client) CreateUpdateDeployment(name, namespace string, conf *config.Config) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	containers := []corev1.Container{}
	for _, container := range conf.Containers {
		containers = append(containers, corev1.Container{
			Name:  container.Name,
			Image: container.Image + ":" + container.Tag,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: container.Port,
				},
			},
		})
	}

	deployment := &appsv1.Deployment{
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
		Spec: appsv1.DeploymentSpec{
			Replicas: &conf.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					appLabelName: name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: map[string]string{
						appLabelName: name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
					// TODO service account
					ServiceAccountName: "",
				},
			},
		},
	}

	_, err = c.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err = c.Clientset.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("updating deployment failed: %s", err)
			}
		} else {
			return fmt.Errorf("creating deployment failed: %s", err)
		}
	}

	return nil
}

// DeleteDeployment delete the specified deployment
func (c *Client) DeleteDeployment(name, namespace string) error {
	if err := c.Clientset.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting deployment failed: %s", err)
	}

	return nil
}
