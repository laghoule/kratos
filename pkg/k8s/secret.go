package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: refactor CreateUpdateSecret to work as the other CreateUpdateObject
// TODO: add a CreateUpdateConfiguration, for storing kratos release configuration

// checkSecretOwnership check if it's safe to create, update or delete the secret
func (c *Client) checkSecretOwnership(name, namespace string) error {
	svc, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting secret failed: %s", err)
	}

	if svc.Labels[depLabelName] == name {
		return nil
	}

	return fmt.Errorf("secret is not owned by kratos")
}

// CreateUpdateSecret a secret to namespace
func (c *Client) CreateUpdateSecret(secret *corev1.Secret, namespace string) error {
	if err := c.checkSecretOwnership(secret.Name, namespace); err != nil {
		return err
	}

	_, err := c.Clientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := c.Clientset.CoreV1().Secrets(namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("updating secret %s failed: %s", secret.Name, err)
			}
		} else {
			return fmt.Errorf("creation of secret %s failed: %s", secret.Name, err)
		}
	}

	return nil
}

// DeleteSecret delete a secret from a namespace
func (c *Client) DeleteSecret(name, namespace string) error {
	if err := c.checkSecretOwnership(name, namespace); err != nil {
		return err
	}

	if err := c.Clientset.CoreV1().Secrets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting secret %s failed: %s", name, err)
	}

	return nil
}

// GetSecret get a secret from a namespace
func (c *Client) GetSecret(name, namespace string) (*corev1.Secret, error) {
	secret, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting secret %s failed: %s", name, err)
	}

	return secret, nil
}
