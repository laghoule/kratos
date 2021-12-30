package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/imdario/mergo"
)

const (
	clusterIssuerAnnotation = "cert-manager.io/cluster-issuer"
	sslRedirectAnnotation   = "nginx.ingress.kubernetes.io/ssl-redirect"
)

// Ingress contain the kubernetes clientset and configuration of the release
type Ingress struct {
	Clientset kubernetes.Interface
	*config.Config
}

// checkOwnership check if it's safe to create, update or delete the ingress
func (i *Ingress) checkOwnership(name, namespace string) error {
	ing, err := i.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting ingress failed: %s", err)
	}

	// managed by kratos
	if err := checkKratosManaged(ing.Labels); err == nil {
		if ing.Labels[DepLabelName] == name {
			return nil
		}
	}

	return fmt.Errorf("ingress is not managed by kratos")
}

// CreateUpdate create or update an ingress
func (i *Ingress) CreateUpdate(name, namespace string) error {
	if err := i.checkOwnership(name, namespace); err != nil {
		return err
	}

	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	// merge common & ingress labels
	if err := mergo.Map(&i.Deployment.Ingress.Labels, i.Common.Labels); err != nil {
		return fmt.Errorf("merging ingress labels failed: %s", err)
	}

	// merge kratosLabels & ingress labels
	if err := mergo.Map(&i.Deployment.Ingress.Labels, map[string]string(kratosLabel)); err != nil {
		return fmt.Errorf("merging ingress labels failed: %s", err)
	}

	// merge common & ingress annotations
	if err := mergo.Map(&i.Deployment.Ingress.Annotations, i.Common.Annotations); err != nil {
		return fmt.Errorf("merging ingress annotations failed: %s", err)
	}

	sslAnnotations := map[string]string{
		clusterIssuerAnnotation: i.Deployment.Ingress.ClusterIssuer,
		sslRedirectAnnotation:   "true",
	}

	// merge ingress annotations & sslAnnotations
	if err := mergo.Map(&i.Deployment.Ingress.Annotations, sslAnnotations); err != nil {
		return fmt.Errorf("merging ingress annotations failed: %s", err)
	}

	ingressTLS := []netv1.IngressTLS{}
	ingressRules := []netv1.IngressRule{}
	pathType := netv1.PathTypePrefix

	for _, hostname := range i.Deployment.Ingress.Hostnames {
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
										Number: i.Deployment.Port,
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
				i.Deployment.Ingress.Labels,
				labels.Set{
					DepLabelName: name,
				},
			),
			Annotations: i.Deployment.Ingress.Annotations,
		},
		Spec: netv1.IngressSpec{
			IngressClassName: &i.Deployment.Ingress.IngressClass,
			TLS:              ingressTLS,
			Rules:            ingressRules,
		},
	}

	_, err = i.Clientset.NetworkingV1().Ingresses(namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := i.Clientset.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("updating ingress failed: %s", err)
			}
		} else {
			return fmt.Errorf("creating ingress failed: %s", err)
		}
	}

	return nil
}

// Delete specified ingress
func (i *Ingress) Delete(name, namespace string) error {
	if err := i.checkOwnership(name, namespace); err != nil {
		return err
	}

	err := i.Clientset.NetworkingV1().Ingresses(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deleting ingress failed: %s", err)
	}

	return nil
}

// IsIngressClassExist check if an ingress class object exist
func (i *Ingress) IsIngressClassExist(name string) error {
	_, err := i.Clientset.NetworkingV1().IngressClasses().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("ingressClass %s not found", name)
	}

	return nil
}
