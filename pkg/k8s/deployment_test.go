package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kratosLabelName  = "app.kubernetes.io/managed-by"
	kratosLabelValue = "kratos"
	name             = "name"
	namespace        = "default"
	image            = "test"
	tagLatest        = "latest"
	tagV1            = "v1.0.0"
	appLabelName     = "app"
	containerPort    = 80
)

var (
	replicas   int32 = 1
	deployment       = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				kratosLabelName: kratosLabelValue,
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
									ContainerPort: containerPort,
								},
							},
						},
					},
				},
			},
		},
	}
)

// TestListDeployment test list of no deployment
func TestListNoDeployment(t *testing.T) {
	client := newTest()

	listDep, err := client.ListDeployments(namespace)
	if err != nil {
		return
	}

	assert.Empty(t, listDep)
}

// TestListDeployment test list of one deployment
func TestListDeployment(t *testing.T) {
	client := newTest()

	if err := client.CreateUpdateDeployment(name, namespace, image, tagLatest, replicas); err != nil {
		t.Error(err)
	}

	listDep, err := client.ListDeployments(namespace)
	if err != nil {
		t.Error(err)
	}

	assert.NotEmpty(t, listDep)
}

// TestCreateDeployment test creation of deployment
func TestCreateDeployment(t *testing.T) {
	client := newTest()

	if err := client.CreateUpdateDeployment(name, namespace, image, tagLatest, replicas); err != nil {
		t.Error(err)
	}

	dep, err := client.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, deployment, dep)
}

// TestUpdateDeployment test update deployment of tag from latest to v1.0.0
func TestUpdateDeployment(t *testing.T) {
	client := newTest()

	if err := client.CreateUpdateDeployment(name, namespace, image, tagLatest, replicas); err != nil {
		t.Error(err)
	}

	if err := client.CreateUpdateDeployment(name, namespace, image, tagV1, replicas); err != nil {
		t.Error(err)
	}

	dep, err := client.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, image+":"+tagV1, dep.Spec.Template.Spec.Containers[0].Image)
}
