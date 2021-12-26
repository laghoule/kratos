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

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
)

// checkDeploymentOwnership check if it's safe to create, update or delete the deployment
func (c *Client) checkDeploymentOwnership(name, namespace string) error {
	svc, err := c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting deployment failed: %s", err)
	}

	// TODO: Should also check for config.DeployLabel

	if svc.Labels[DepLabelName] == name {
		return nil
	}

	return fmt.Errorf("deployment is not owned by kratos")
}

// ListDeployments list deployments
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
	if err := c.checkDeploymentOwnership(name, namespace); err != nil {
		return err
	}

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
	volumesMount := []corev1.VolumeMount{}
	volumes := []corev1.Volume{}

	for _, container := range conf.Deployment.Containers {
		volumesMount, volumes = getVolumesConfForContainer(name, &container, conf)

		containers = append(containers, corev1.Container{
			Name:  container.Name,
			Image: container.Image + ":" + container.Tag,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: conf.Deployment.Port,
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
			Annotations: conf.Deployment.Annotations,
			Labels: labels.Merge(
				conf.Deployment.Labels,
				labels.Set{
					DepLabelName: name,
				},
			),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &conf.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					DepLabelName: name,
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
							DepLabelName: name,
						},
					),
				},
				Spec: corev1.PodSpec{
					Containers:                   containers,
					AutomountServiceAccountToken: boolPTR(false),
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: boolPTR(true),
					},
					Volumes: volumes,
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
	if err := c.checkDeploymentOwnership(name, namespace); err != nil {
		return err
	}

	if err := c.Clientset.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting deployment failed: %s", err)
	}

	return nil
}
