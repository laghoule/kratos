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

var (
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				kratosLabelName: kratosLabelValue,
				appLabelName:    name,
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
				appLabelName: name,
			},
		},
	}
)

// TestCreateUpdateDeployment test creation of deployment
func TestCreateUpdateService(t *testing.T) {
	client := testNew()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
	}

	if err := client.CreateUpdateService(name, namespace, conf); err != nil {
		t.Error(err)
	}

	svc, err := client.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, service, svc)
}

// TestCreateDeployment test creation of deployment
func TestUpdateService(t *testing.T) {
	client := testNew()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
	}

	if err := client.CreateUpdateService(name, namespace, conf); err != nil {
		t.Error(err)
	}

	conf.Service.Port = 443

	if err := client.updateService(name, namespace, conf); err != nil {
		t.Error(err)
	}

	svc, err := client.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, int32(443), svc.Spec.Ports[0].Port)
}

func TestDeleteService(t *testing.T) {
	client := testNew()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
	}

	if err := client.CreateUpdateService(name, namespace, conf); err != nil {
		t.Error(err)
	}

	svc, err := client.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}

	assert.NotEmpty(t, svc)

	if err := client.DeleteService(name, namespace); err != nil {
		t.Error(err)
	}

	svc, err = client.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Error(err)
	}

	assert.True(t, errors.IsNotFound(err))
}
