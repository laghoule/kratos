package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Secrets is the interface for secrets
type Secrets interface {
	CreateUpdate(string, string) error
	Delete(string, string) error
	DeleteConfig(string, string) error
	Get(string, string) (*corev1.Secret, error)
	List(string) ([]corev1.Secret, error)
	SaveConfig(string, string, string, string) error
}

// secrets contain the kubernetes clientset and configuration of the release
type secrets struct {
	Clientset kubernetes.Interface
	*config.Config
}

// checkOwnership check if it's safe to create, update or delete the secret
func (s *secrets) checkOwnership(name, namespace string) error {
	secret, err := s.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting secret failed: %s", err)
	}

	// managed by kratos
	if err := checkKratosManaged(secret.Labels); err == nil {
		if secret.Labels[SecretLabelName] == name {
			return nil
		}
	}

	return fmt.Errorf("secret is not managed by kratos")
}

// SaveConfig save kratos release configuration
func (s *secrets) SaveConfig(name, namespace, key, value string) error {
	if err := s.checkOwnership(name, namespace); err != nil {
		return err
	}

	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return nil
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					SecretLabelName: name,
				},
			),
		},
		StringData: map[string]string{
			key: value,
		},
		Type: corev1.SecretTypeOpaque,
	}

	if err := s.createUpdate(secret); err != nil {
		return err
	}

	return nil
}

// DeleteConfig delete kratos release configuration
func (s *secrets) DeleteConfig(name, namespace string) error {
	if err := s.delete(name, namespace); err != nil {
		return err
	}

	return nil
}

// createUpdate create or update a secret
func (s *secrets) createUpdate(secret *corev1.Secret) error {
	if err := s.checkOwnership(secret.Name, secret.Namespace); err != nil {
		return err
	}

	_, err := s.Clientset.CoreV1().Secrets(secret.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := s.Clientset.CoreV1().Secrets(secret.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("updating secret %s failed: %s", secret.Name, err)
			}
		} else {
			return fmt.Errorf("creation of secret %s failed: %s", secret.Name, err)
		}
	}

	return nil
}

// CreateUpdate create or update a secrets with value provided in conf
func (s *secrets) CreateUpdate(name, namespace string) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return nil
	}

	if s.Secrets != nil {
		// merge labels
		if err := mergeStringMaps(&s.Secrets.Labels, s.Common.Labels, kratosLabel); err != nil {
			return fmt.Errorf("merging secret labels failed: %s", err)
		}

		// merge annotations
		if err := mergeStringMaps(&s.Secrets.Annotations, s.Common.Annotations); err != nil {
			return fmt.Errorf("merging secrets annotations failed: %s", err)
		}
	} else {
		// merge kratosLabels & common labels
		s.Secrets = &config.Secrets{
			Labels: map[string]string(kratosLabel),
		}
		if err := mergeStringMaps(&s.Secrets.Labels, s.Common.Labels); err != nil {
			return fmt.Errorf("merging secrets labels failed: %s", err)
		}
	}

	for _, file := range s.Secrets.Files {
		secretName := name + "-" + file.Name

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
				Labels: labels.Merge(
					s.Secrets.Labels,
					labels.Set{
						SecretLabelName: secretName,
					},
				),
				Annotations: s.Common.Annotations,
			},
			StringData: map[string]string{
				file.Name: file.Data,
			},
			Type: corev1.SecretTypeOpaque,
		}

		if err := s.createUpdate(secret); err != nil {
			return err
		}
	}

	return nil
}

// Delete the secrets contained in conf for the specified namespace
func (s *secrets) Delete(name, namespace string) error {
	for _, file := range s.Secrets.Files {
		if err := s.delete(name+"-"+file.Name, namespace); err != nil {
			return err
		}
	}

	return nil
}

// delete a secret from a namespace
func (s *secrets) delete(name, namespace string) error {
	if err := s.checkOwnership(name, namespace); err != nil {
		return err
	}

	if err := s.Clientset.CoreV1().Secrets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting secret %s failed: %s", name, err)
	}

	return nil
}

// Get a secret from a namespace
func (s *secrets) Get(name, namespace string) (*corev1.Secret, error) {
	secret, err := s.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting secret %s failed: %s", name, err)
	}

	return secret, nil
}

// List the secret in the specified namespace
func (s *secrets) List(namespace string) ([]corev1.Secret, error) {
	list, err := s.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: config.ManagedLabel,
	})
	if err != nil {
		return nil, fmt.Errorf("getting secrets list failed: %s", err)
	}

	return list.Items, nil
}
