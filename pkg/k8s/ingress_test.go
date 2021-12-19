package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	netv1 "k8s.io/api/networking/v1"
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
				DepLabelName:    name,
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

func createNotKratosIngress(c *Client) error {
	ing := createIngress()
	ing.Labels = nil

	if _, err := c.Clientset.NetworkingV1().Ingresses(namespace).Create(context.Background(), ing, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

// TestCreateUpdateIngressNotOwnedByKratos test update of an ingress not owned by kratos
func TestCreateUpdateIngressNotOwnedByKratos(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := createNotKratosIngress(c); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateIngress(name, namespace, conf); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "ingress is not owned by kratos")
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
	if err := c.CreateUpdateIngress(name, namespace, conf); err != nil {
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
	conf.Deployment.Ingress.Hostnames[0] = "www.example.com"
	if err = c.CreateUpdateIngress(name, namespace, conf); err != nil {
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

// TestDeleteIngressNotOwnedByKratos test removing of ingress not owned by kratos
func TestDeleteIngressNotOwnedByKratos(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := createNotKratosIngress(c); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateIngress(name, namespace, conf); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "ingress is not owned by kratos")
	}
}

// TestDeleteIngress test removing of ingress
func TestDeleteIngress(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateIngress(name, namespace, conf); err != nil {
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

	list, err := c.Clientset.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Empty(t, list)
}

// TestIsIngressClassExist test if an ingressClass exist
func TestIsIngressClassExist(t *testing.T) {
	c := new()
	class := createIngressClass()

	if _, err := c.Clientset.NetworkingV1().IngressClasses().Create(context.Background(), class, metav1.CreateOptions{}); err != nil {
		t.Error(err)
		return
	}

	err := c.IsIngressClassExist(name)
	assert.NoError(t, err)
}
