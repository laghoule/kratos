package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// Service is the interface for service
type Service interface {
	CreateUpdate(string, string) error
	Delete(string, string) error
	List(string) ([]corev1.Service, error)
}

// service contain the kubernetes clientset and configuration of the release
type service struct {
	Clientset kubernetes.Interface
	*config.Config
}

// checkOwnership check if it's safe to create, update or delete the service
func (s *service) checkOwnership(name, namespace string) error {
	svc, err := s.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting service failed: %s", err)
	}

	// managed by kratos
	if err := checkKratosManaged(svc.Labels); err == nil {
		if svc.Labels[DepLabelName] == name {
			return nil
		}
	}

	return fmt.Errorf("service is not managed by kratos")
}

// CreateUpdate create or update a service
func (s *service) CreateUpdate(name, namespace string) error {
	if err := s.checkOwnership(name, namespace); err != nil {
		return err
	}

	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return nil
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					DepLabelName: name,
				},
			),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: name,
					Port: s.Deployment.Port,
					TargetPort: intstr.IntOrString{
						IntVal: s.Deployment.Port,
					},
				},
			},
			Selector: map[string]string{
				DepLabelName: name,
			},
		},
	}

	_, err = s.Clientset.CoreV1().Services(namespace).Create(context.Background(), svc, metav1.CreateOptions{})
	if err != nil {
		// if service exist, we call update
		if errors.IsAlreadyExists(err) {
			if err := s.update(name, namespace); err != nil {
				return fmt.Errorf("updating service failed: %s", err)
			}
		} else {
			return fmt.Errorf("creating service failed: %s", err)
		}
	}

	return nil
}

// update an existing service. Used by CreateUpdateService.
func (s *service) update(name, namespace string) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return nil
	}

	svcLabels := map[string]string{}
	if err := mergeStringMaps(&svcLabels, s.Common.Labels, kratosLabel); err != nil {
		return err
	}

	svcInfo, err := s.Clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("getting service informations failed: %s", err)
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				svcLabels,
				labels.Set{
					DepLabelName: name,
				},
			),
			Annotations:     s.Common.Annotations,
			ResourceVersion: svcInfo.ResourceVersion,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: svcInfo.Spec.ClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: name,
					Port: s.Deployment.Port,
					TargetPort: intstr.IntOrString{
						IntVal: s.Deployment.Port,
					},
				},
			},
			Selector: map[string]string{
				DepLabelName: name,
			},
		},
	}

	_, err = s.Clientset.CoreV1().Services(namespace).Update(context.Background(), svc, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("updating service failed: %s", err)
	}

	return nil
}

// Delete the specified service
func (s *service) Delete(name, namespace string) error {
	if err := s.checkOwnership(name, namespace); err != nil {
		return err
	}

	if err := s.Clientset.CoreV1().Services(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("delete service failed: %s", err)
	}

	return nil
}

// List services of the specified namespace
func (s *service) List(namespace string) ([]corev1.Service, error) {
	list, err := s.Clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: config.ManagedLabel,
	})
	if err != nil {
		return nil, fmt.Errorf("getting services list failed: %s", err)
	}

	return list.Items, nil
}
