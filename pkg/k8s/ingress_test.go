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

// createIngressTLS return a list of IngressTLS object
func createInressTLS() []netv1.IngressTLS {
	return []netv1.IngressTLS{
		{
			Hosts:      []string{hostname},
			SecretName: hostname + "-tls",
		},
	}
}

// createIngressRules return an ingressRule object
func createIngressRules() []netv1.IngressRule {
	var pathType = netv1.PathTypePrefix
	return []netv1.IngressRule{
		{
			Host: hostname,
			IngressRuleValue: netv1.IngressRuleValue{
				HTTP: &netv1.HTTPIngressRuleValue{
					Paths: []netv1.HTTPIngressPath{
						{
							Path:     path,
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
}

// createIngress return an ingress object
func createIngress() *netv1.Ingress {
	ingressClass := "nginx"
	return &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				kratosLabelName: kratosLabelValue,
				depLabelName:    name,
				"cloudflare":    "enabled",
				"environment":   "dev",
			},
			Annotations: map[string]string{
				clusterIssuerAnnotation: clusterIssuer,
				sslRedirectAnnotation:   "true",
				"branch":                "dev",
				"hsts":                  "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: &ingressClass,
			TLS:              createInressTLS(),
			Rules:            createIngressRules(),
		},
	}
}

// createIngress return an ingressClass object
func createIngressClass() *netv1.IngressClass {
	ingressClass := "nginx"
	return &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: netv1.IngressClassSpec{
			Controller: ingressClass,
		},
	}
}

// TestCreateUpdateIngress test the creation and update of an ingress
func TestCreateUpdateIngress(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	// create
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

	assert.Equal(t, createIngress(), ing)

	// update
	conf.Ingress.Hostnames[0] = "www.example.com"
	err = c.CreateUpdateIngress(name, namespace, conf)
	if err != nil {
		t.Error(err)
		return
	}

	ing, err = c.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "www.example.com", ing.Spec.Rules[0].Host)
}

// TestDeleteIngress test removing of ingress
func TestDeleteIngress(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
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

// TestIsIngressClassExist test if an ingressClass exist
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
