package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	goodConfig = "../config/testdata/goodConfig.yml"
)

var (
	ingressClass = "nginx"

	pathType = netv1.PathTypePrefix

	ingressTLS = []netv1.IngressTLS{
		{
			Hosts:      []string{hostname},
			SecretName: hostname + "-tls",
		},
	}

	ingressRules = []netv1.IngressRule{
		{
			Host: hostname,
			IngressRuleValue: netv1.IngressRuleValue{
				HTTP: &netv1.HTTPIngressRuleValue{
					Paths: []netv1.HTTPIngressPath{
						{
							Path:     "/",
							PathType: &pathType,
							Backend: netv1.IngressBackend{
								Service: &netv1.IngressServiceBackend{
									Name: name,
									Port: netv1.ServiceBackendPort{
										Number: containerHTTP,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ingress = &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				kratosLabelName: kratosLabelValue,
				appLabelName:    name,
			},
			Annotations: map[string]string{
				clusterIssuerAnnotation: clusterIssuer,
				sslRedirectAnnotation:   "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: &ingressClass,
			TLS:              ingressTLS,
			Rules:            ingressRules,
		},
	}
)

func createIngressClass() *netv1.IngressClass {
	return &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: netv1.IngressClassSpec{
			Controller: ingressClass,
		},
	}
}

func TestCreateIngress(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
		return
	}

	err := c.CreateUpdateIngress(name, namespace, conf)
	if err != nil {
		t.Error(err)
		return
	}

	ing, err := c.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, ingress, ing)
}

func TestUpdateIngress(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
		return
	}

	err := c.CreateUpdateIngress(name, namespace, conf)
	if err != nil {
		t.Error(err)
		return
	}

	err = c.CreateUpdateIngress(name, namespace, conf)
	if err != nil {
		t.Error(err)
		return
	}

	ing, err := c.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "example.com", ing.Spec.Rules[0].Host)
}

func TestDeleteIngress(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(goodConfig); err != nil {
		t.Error(err)
		return
	}

	err := c.CreateUpdateIngress(name, namespace, conf)
	if err != nil {
		t.Error(err)
		return
	}

	ing, err := c.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, ing)

	if err := c.DeleteIngress(name, namespace); err != nil {
		t.Error(err)
		return
	}

	ing, err = c.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Error(err)
		return
	}

	assert.True(t, errors.IsNotFound(err))
}

func TestIsIngressClassExist(t *testing.T) {
	c := new()
	conf := &config.Config{
		Ingress: &config.Ingress{},
	}

	conf.ClusterIssuer = clusterIssuer
	class := createIngressClass()

	_, err := c.Clientset.NetworkingV1().IngressClasses().Create(context.Background(), class, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	list, err := c.Clientset.NetworkingV1().IngressClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	found := c.IsIngressClassExist(name)

	assert.True(t, found)
}
