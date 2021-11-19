package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createDeployment() *appsv1.Deployment {
	var replicas int32 = 1
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				kratosLabelName: kratosLabelValue,
				appLabelName:    name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
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
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image + ":" + tagLatest,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: containerHTTP,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{},
								Limits:   corev1.ResourceList{},
							},
						},
					},
				},
			},
		},
	}
}

// TestListNoDeployment test list of no deployment
func TestListNoDeployment(t *testing.T) {
	c := new()

	list, err := c.ListDeployments(namespace)
	if err != nil {
		return
	}

	assert.Empty(t, list)
}

// TestListDeployment test list of one deployment
func TestListDeployment(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateDeployment(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	list, err := c.ListDeployments(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, list)
}

// TestCreateDeployment test creation of deployment
func TestCreateDeployment(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateDeployment(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	dep, err := c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, createDeployment(), dep)
}

// TestUpdateDeployment test update deployment of tag from latest to v1.0.0
func TestUpdateDeployment(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateDeployment(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	conf.Containers[0].Tag = "v1.0.0"

	if err := c.CreateUpdateDeployment(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	dep, err := c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, image+":"+tagV1, dep.Spec.Template.Spec.Containers[0].Image)
}

func TestDeleteDeployment(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateDeployment(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	dep, err := c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, dep)

	if err := c.DeleteDeployment(name, namespace); err != nil {
		t.Error(err)
		return
	}

	dep, err = c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Error(err)
		return
	}

	assert.True(t, errors.IsNotFound(err))
}
