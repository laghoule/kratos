package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/imdario/mergo"
)

const (
	clusterIssuerAnnotation = "cert-manager.io/cluster-issuer"
	sslRedirectAnnotation   = "nginx.ingress.kubernetes.io/ssl-redirect"
)

// checkIngressOwnership check if it's safe to create, update or delete the ingress
func (c *Client) checkIngressOwnership(name, namespace string) error {
	svc, err := c.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting ingress failed: %s", err)
	}

	// TODO: Should also check for config.DeployLabel

	if svc.Labels[DepLabelName] == name {
		return nil
	}

	return fmt.Errorf("ingress is not owned by kratos")
}

// CreateUpdateIngress create or update an ingress
func (c *Client) CreateUpdateIngress(name, namespace string, conf *config.Config) error {
	if err := c.checkIngressOwnership(name, namespace); err != nil {
		return err
	}

	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	// merge common & ingress labels
	if err := mergo.Map(&conf.Deployment.Ingress.Labels, conf.Common.Labels); err != nil {
		return fmt.Errorf("merging ingress labels failed: %s", err)
	}

	// merge kratosLabels & ingress labels
	if err := mergo.Map(&conf.Deployment.Ingress.Labels, map[string]string(kratosLabel)); err != nil {
		return fmt.Errorf("merging ingress labels failed: %s", err)
	}

	// merge common & ingress annotations
	if err := mergo.Map(&conf.Deployment.Ingress.Annotations, conf.Common.Annotations); err != nil {
		return fmt.Errorf("merging ingress annotations failed: %s", err)
	}

	sslAnnotations := map[string]string{
		clusterIssuerAnnotation: conf.Deployment.Ingress.ClusterIssuer,
		sslRedirectAnnotation:   "true",
	}

	// merge ingress annotations & sslAnnotations
	if err := mergo.Map(&conf.Deployment.Ingress.Annotations, sslAnnotations); err != nil {
		return fmt.Errorf("merging ingress annotations failed: %s", err)
	}

	ingressTLS := []netv1.IngressTLS{}
	ingressRules := []netv1.IngressRule{}
	pathType := netv1.PathTypePrefix

	for _, hostname := range conf.Deployment.Ingress.Hostnames {
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
										Number: conf.Deployment.Port,
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
				conf.Deployment.Ingress.Labels,
				labels.Set{
					DepLabelName: name,
				},
			),
			Annotations: conf.Deployment.Ingress.Annotations,
		},
		Spec: netv1.IngressSpec{
			IngressClassName: &conf.Deployment.Ingress.IngressClass,
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
	if err := c.checkIngressOwnership(name, namespace); err != nil {
		return err
	}

	err := c.Clientset.NetworkingV1().Ingresses(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deleting ingress failed: %s", err)
	}

	return nil
}

// IsIngressClassExist check if an ingress class object exist
func (c *Client) IsIngressClassExist(name string) error {
	_, err := c.Clientset.NetworkingV1().IngressClasses().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("ingressClass %s not found", name)
	}

	return nil
}
