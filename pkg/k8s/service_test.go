package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// createService return a service object
func createService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				kratosLabelName: kratosLabelValue,
				depLabelName:    name,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: name,
					Port: containerHTTP,
					TargetPort: intstr.IntOrString{
						IntVal: containerHTTP,
					},
				},
			},
			Selector: map[string]string{
				depLabelName: name,
			},
		},
	}
}

// TestCreateUpdateServiceNotOwnedByKratos test update of a service not owned by kratos
func TestCreateUpdateServiceNotOwnedByKratos(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	extSvc := createService()
	extSvc.Labels = nil

	_, err := c.Clientset.CoreV1().Services(namespace).Create(context.Background(), extSvc, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	// create && fail
	if err := c.CreateUpdateService(name, namespace, conf); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "service is not owned by kratos")
	}
}

// TestCreateUpdateDeployment test creation and update of a service
func TestCreateUpdateService(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	// create
	if err := c.CreateUpdateService(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	svc, err := c.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, createService(), svc)

	// update
	conf.Deployment.Port = 443
	if err := c.CreateUpdateService(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	svc, err = c.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, int32(443), svc.Spec.Ports[0].Port)
}

// TestDeleteService test delete of a service not owned by kratos
func TestDeleteServiceNotOwnedByKratos(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	extSvc := createService()
	extSvc.Labels = nil

	_, err := c.Clientset.CoreV1().Services(namespace).Create(context.Background(), extSvc, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if err := c.DeleteService(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "service is not owned by kratos")
	}
}

// TestDeleteService test delete of a service
func TestDeleteService(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateService(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	svc, err := c.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, svc)

	if err := c.DeleteService(name, namespace); err != nil {
		t.Error(err)
		return
	}

	_, err = c.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Error(err)
		return
	}

	assert.True(t, errors.IsNotFound(err))
}
