package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/common"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	clusterIssuerAnnotation = "cert-manager.io/cluster-issuer"
	sslRedirectAnnotation   = "nginx.ingress.kubernetes.io/ssl-redirect"
)

// CreateUpdateIngress create or update an ingress
func (c *Client) CreateUpdateIngress(name, namespace, ingressClass, clusterIssuer string, hostnames []string, port int32) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(common.DeployLabel)
	if err != nil {
		return nil
	}

	ingressTLS := []netv1.IngressTLS{}
	ingressRules := []netv1.IngressRule{}
	pathType := netv1.PathTypePrefix

	for _, hostname := range hostnames {
		ingressTLS = append(ingressTLS, netv1.IngressTLS{
			Hosts:      []string{hostname},
			SecretName: hostname + "-tls",
		})

		ingressRules = append(ingressRules, netv1.IngressRule{
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
										Number: port,
									},
								},
							},
						},
					},
				},
			},
		})
	}

	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					appLabelName: name,
				},
			),
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

	_, err = c.Clientset.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := c.Clientset.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("updating ingress failed: %s", err)
			}
		} else {
			return fmt.Errorf("creating ingress failed: %s", err)
		}
	}

	return nil
}

// DeleteIngress delete specified ingress
func (c *Client) DeleteIngress(name, namespace string) error {
	err := c.Clientset.NetworkingV1().Ingresses(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deleting ingress failed: %s", err)
	}

	return nil
}
