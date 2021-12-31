package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func newIngress() (*Ingress, error) {
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		return nil, err
	}

	return &Ingress{
		Clientset: fake.NewSimpleClientset(),
		Config:    conf,
	}, nil
}

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

func createNotKratosIngress(i *Ingress) error {
	ing := createIngress()
	ing.Labels = nil

	if _, err := i.Clientset.NetworkingV1().Ingresses(namespace).Create(context.Background(), ing, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

// TestCreateUpdateIngressNotOwnedByKratos test update of an ingress not owned by kratos
func TestCreateUpdateIngressNotOwnedByKratos(t *testing.T) {
	i, err := newIngress()
	if err != nil {
		t.Error(err)
		return
	}

	if err := createNotKratosIngress(i); err != nil {
		t.Error(err)
		return
	}

	if err := i.CreateUpdate(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "ingress is not managed by kratos")
	}
}

// TestCreateUpdate test the creation and update of an ingress
func TestCreateUpdateIngress(t *testing.T) {
	i, err := newIngress()
	if err != nil {
		t.Error(err)
		return
	}

	// create
	if err := i.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	ing, err := i.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, createIngress(), ing)

	// update
	i.Deployment.Ingress.Hostnames[0] = "www.example.com"
	if err = i.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	ing, err = i.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "www.example.com", ing.Spec.Rules[0].Host)
}

// TestDeleteIngressNotOwnedByKratos test removing of ingress not owned by kratos
func TestDeleteIngressNotOwnedByKratos(t *testing.T) {
	i, err := newIngress()
	if err != nil {
		t.Error(err)
		return
	}

	if err := createNotKratosIngress(i); err != nil {
		t.Error(err)
		return
	}

	if err := i.CreateUpdate(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "ingress is not managed by kratos")
	}
}

// TestDeleteIngress test removing of ingress
func TestDeleteIngress(t *testing.T) {
	i, err := newIngress()
	if err != nil {
		t.Error(err)
		return
	}

	if err := i.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	ing, err := i.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, ing)

	if err := i.Delete(name, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := i.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Empty(t, list)
}

// TestIsIngressClassExist test if an ingressClass exist
func TestIsIngressClassExist(t *testing.T) {
	class := createIngressClass()

	i, err := newIngress()
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := i.Clientset.NetworkingV1().IngressClasses().Create(context.Background(), class, metav1.CreateOptions{}); err != nil {
		t.Error(err)
		return
	}

	err = i.IsIngressClassExist(name)
	assert.NoError(t, err)
}
