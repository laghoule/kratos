package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
)

const (
	resCPU    = "cpu"
	resMemory = "memory"
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

// formatResources format the resource from container configurations
func formatResources(container config.Container) corev1.ResourceRequirements {
	req := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{},
		Limits:   corev1.ResourceList{},
	}

	// requests
	if container.Resources.Request.CPU != "" {
		req.Requests[resCPU] = resource.MustParse(container.Resources.Request.CPU)
	}
	if container.Resources.Request.Memory != "" {
		req.Requests[resMemory] = resource.MustParse(container.Resources.Request.Memory)
	}

	// limits
	if container.Resources.Limits.CPU != "" {
		req.Limits[resCPU] = resource.MustParse(container.Resources.Limits.CPU)
	}
	if container.Resources.Limits.Memory != "" {
		req.Requests[resMemory] = resource.MustParse(container.Resources.Limits.Memory)
	}

	return req
}

// CreateUpdateDeployment create or update a deployment
func (c *Client) CreateUpdateDeployment(name, namespace string, conf *config.Config) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return fmt.Errorf("converting label failed: %s", err)
	}

	// merge common & deployments labels
	if err := mergo.Map(&conf.Deployment.Labels, conf.Common.Labels); err != nil {
		return fmt.Errorf("merging deployment labels failed: %s", err)
	}

	// deep copy of map for podLabels
	podLabels := map[string]string{}
	if err := copier.Copy(&podLabels, &conf.Deployment.Labels); err != nil {
		return fmt.Errorf("copying deployment labels values failed: %s", err)
	}

	// merge kratosLabels & deployment labels
	if err := mergo.Map(&conf.Deployment.Labels, map[string]string(kratosLabel)); err != nil {
		return fmt.Errorf("merging deployment labels failed: %s", err)
	}

	// merge common & deployments annotations
	if err := mergo.Map(&conf.Deployment.Annotations, conf.Common.Annotations); err != nil {
		return fmt.Errorf("merging deployment annotations failed: %s", err)
	}

	containers := []corev1.Container{}

	for _, container := range conf.Containers {
		containers = append(containers, corev1.Container{
			Name:  container.Name,
			Image: container.Image + ":" + container.Tag,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: conf.Deployment.Port,
				},
			},
			Resources: formatResources(container),
		})
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: conf.Deployment.Annotations,
			Labels: labels.Merge(
				conf.Deployment.Labels,
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
					Name:        name,
					Namespace:   namespace,
					Annotations: conf.Deployment.Annotations,
					Labels: labels.Merge(
						podLabels,
						labels.Set{
							appLabelName: name,
						},
					),
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
