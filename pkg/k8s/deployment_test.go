package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func createDeployment() *appsv1.Deployment {
	var replicas int32 = 1
	depLabels := map[string]string{
		kratosLabelName: kratosLabelValue,
		depLabelName:    name,
		"environment":   environment,
		"app":           name,
	}
	podLabels := map[string]string{
		depLabelName:  name,
		"environment": environment,
		"app":         name,
	}
	annotations := map[string]string{
		"branch":   environment,
		"revision": "22",
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      depLabels,
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					depLabelName: name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        name,
					Namespace:   namespace,
					Labels:      podLabels,
					Annotations: annotations,
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
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("25m"),
									corev1.ResourceMemory: resource.MustParse("32Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("50m"),
									corev1.ResourceMemory: resource.MustParse("64Mi"),
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/isLive",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 80,
										},
									},
								},
								InitialDelaySeconds: 3,
								PeriodSeconds:       3,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/isReady",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 80,
										},
									},
								},
								InitialDelaySeconds: 3,
								PeriodSeconds:       3,
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

	if err := conf.Load(deploymentConfig); err != nil {
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

// TestCreateUpdateDeploymentNotOwnedByKratos test update of deployment not owned by kratos
func TestCreateUpdateDeploymentNotOwnedByKratos(t *testing.T) {
	c := new()
	conf := &config.Config{}

	dep := createDeployment()
	dep.Labels = nil

	_, err := c.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), dep, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateDeployment(name, namespace, conf); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "deployment is not owned by kratos")
	}
}

// TestCreateUpdateDeployment test creation of deployment
func TestCreateUpdateDeployment(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	// create
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

	// update
	conf.Containers[0].Tag = "v1.0.0"
	if err := c.CreateUpdateDeployment(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	dep, err = c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, image+":"+tagV1, dep.Spec.Template.Spec.Containers[0].Image)
}

// TestDeleteDeploymentNotOwnedByKratos test delete of deployment not owned by kratos
func TestDeleteDeploymentNotOwnedByKratos(t *testing.T) {
	c := new()

	dep := createDeployment()
	dep.Labels = nil

	_, err := c.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), dep, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if err := c.DeleteDeployment(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "deployment is not owned by kratos")
	}
}

// TestDeleteDeployment test delete of deployment
func TestDeleteDeployment(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
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

	_, err = c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Error(err)
		return
	}

	assert.True(t, errors.IsNotFound(err))
}
