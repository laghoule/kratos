package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/imdario/mergo"
)

// checkSecretOwnership check if it's safe to create, update or delete the secret
func (c *Client) checkSecretOwnership(name, namespace string) error {
	svc, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting secret failed: %s", err)
	}

	if svc.Labels[SecretLabelName] == name {
		return nil
	}

	return fmt.Errorf("secret is not owned by kratos")
}

// SaveConfig save kratos release configuration
func (c *Client) SaveConfig(name, namespace, key, value string, conf *config.Config) error {
	if err := c.checkSecretOwnership(name, namespace); err != nil {
		return err
	}

	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
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

	if err := c.createUpdateSecret(secret); err != nil {
		return err
	}

	return nil
}

// createUpdateSecret create or update a secret
func (c *Client) createUpdateSecret(secret *corev1.Secret) error {
	_, err := c.Clientset.CoreV1().Secrets(secret.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err := c.Clientset.CoreV1().Secrets(secret.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("updating secret %s failed: %s", secret.Name, err)
			}
		} else {
			return fmt.Errorf("creation of secret %s failed: %s", secret.Name, err)
		}
	}

	return nil
}

// CreateUpdateSecrets create or update a secrets with value provided in conf
func (c *Client) CreateUpdateSecrets(namespace string, conf *config.Config) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	if conf.Secrets != nil {
		// merge common & secrets labels
		if err := mergo.Map(&conf.Secrets.Labels, conf.Common.Labels); err != nil {
			return fmt.Errorf("merging secrets labels failed: %s", err)
		}

		// merge kratosLabels & secrets labels
		if err := mergo.Map(&conf.Secrets.Labels, map[string]string(kratosLabel)); err != nil {
			return fmt.Errorf("merging ingress labels failed: %s", err)
		}

		// merge common & ingress annotations
		if err := mergo.Map(&conf.Secrets.Annotations, conf.Common.Annotations); err != nil {
			return fmt.Errorf("merging secret annotations failed: %s", err)
		}
	} else {
		// merge kratosLabels & common labels
		conf.Secrets = &config.Secrets{
			Labels: map[string]string(kratosLabel),
		}
		if err := mergo.Map(&conf.Secrets.Labels, conf.Common.Labels); err != nil {
			return fmt.Errorf("merging secrets labels failed: %s", err)
		}
	}

	for _, file := range conf.Secrets.Files {
		if err := c.checkSecretOwnership(file.Name, namespace); err != nil {
			return err
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      file.Name,
				Namespace: namespace,
				Labels: labels.Merge(
					conf.Secrets.Labels,
					labels.Set{
						SecretLabelName: file.Name,
					},
				),
				Annotations: conf.Common.Annotations,
			},
			StringData: map[string]string{
				file.Name: file.Data,
			},
			Type: corev1.SecretTypeOpaque,
		}

		if err := c.createUpdateSecret(secret); err != nil {
			return err
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
