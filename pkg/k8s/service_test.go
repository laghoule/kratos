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
	"k8s.io/client-go/kubernetes/fake"
)

func newService() (*Service, error) {
	conf := &config.Config{}

	if err := conf.Load(secretConfig); err != nil {
		return nil, err
	}

	return &Service{
		Clientset: fake.NewSimpleClientset(),
		Config:    conf,
	}, nil
}

// createService return a service object
func createService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				kratosLabelName: kratosLabelValue,
				DepLabelName:    name,
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
				DepLabelName: name,
			},
		},
	}
}

// TestCreateUpdateServiceNotOwnedByKratos test update of a service not owned by kratos
func TestCreateUpdateServiceNotOwnedByKratos(t *testing.T) {
	s, err := newService()
	if err != nil {
		t.Error(err)
		return
	}

	if err := s.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	extSvc := createService()
	extSvc.Labels = nil

	_, err = s.Clientset.CoreV1().Services(namespace).Create(context.Background(), extSvc, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	// create && fail
	if err := s.CreateUpdate(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "service is not managed by kratos")
	}
}

// TestCreateUpdateDeployment test creation and update of a service
func TestCreateUpdateService(t *testing.T) {
	s, err := newService()
	if err != nil {
		t.Error(err)
		return
	}

	// create
	if err := s.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	svc, err := s.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, createService(), svc)

	// update
	s.Deployment.Port = 443
	if err := s.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	svc, err = s.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, int32(443), svc.Spec.Ports[0].Port)
}

// TestDeleteService test delete of a service not owned by kratos
func TestDeleteServiceNotOwnedByKratos(t *testing.T) {
	s, err := newService()
	if err != nil {
		t.Error(err)
		return
	}

	extSvc := createService()
	extSvc.Labels = nil

	_, err = s.Clientset.CoreV1().Services(namespace).Create(context.Background(), extSvc, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if err := s.Delete(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "service is not managed by kratos")
	}
}

// TestDeleteService test delete of a service
func TestDeleteService(t *testing.T) {
	s, err := newService()
	if err != nil {
		t.Error(err)
		return
	}

	if err := s.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	svc, err := s.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, svc)

	if err := s.Delete(name, namespace); err != nil {
		t.Error(err)
		return
	}

	_, err = s.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Error(err)
		return
	}

	assert.True(t, errors.IsNotFound(err))
}
