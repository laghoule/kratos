package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	secretData        = "usename: patate\npassword: poil\n"
	secretUpdatedData = "my updated secret data"
	secretFileName    = "credentials.yaml"
)

func newSecret() (*Secrets, error) {
	conf := &config.Config{}

	if err := conf.Load(secretConfig); err != nil {
		return nil, err
	}

	return &Secrets{
		Clientset: fake.NewSimpleClientset(),
		Config:    conf,
	}, nil
}

// createSecret return a secret object
func createSecret() *corev1.Secret {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return nil
	}

	secretName := name + "-" + secretFileName

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					"app":           name,
					"environment":   environment,
					SecretLabelName: secretName,
				}),
			Annotations: map[string]string{
				"branch": environment,
			},
		},
		StringData: map[string]string{
			secretFileName: secretData,
		},
		Type: corev1.SecretTypeOpaque,
	}
}

// createConfigSecret return a kratos release configuration secret object
func createConfigSecret() *corev1.Secret {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return nil
	}

	secretName := name + config.ConfigSuffix

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					SecretLabelName: secretName,
				}),
		},
		StringData: map[string]string{
			config.ConfigKey: secretData,
		},
		Type: corev1.SecretTypeOpaque,
	}
}

func createNotKratosSecret(s *Secrets) error {
	secret := createSecret()
	secret.Labels = nil

	if _, err := s.Clientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

// TestSaveConfig test saving a release configuration
func TestSaveConfig(t *testing.T) {
	s, err := newSecret()
	if err != nil {
		t.Error(err)
		return
	}

	secretName := name + config.ConfigSuffix

	if err := s.SaveConfig(secretName, namespace, config.ConfigKey, secretData); err != nil {
		t.Error(err)
		return
	}

	secret, err := s.Get(secretName, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	expected := createConfigSecret()
	assert.Equal(t, expected, secret)
}

// TestDeleteConfig test deleting a release configuration
func TestDeleteConfig(t *testing.T) {
	s, err := newSecret()
	if err != nil {
		t.Error(err)
		return
	}

	secretName := name + config.ConfigSuffix

	if err := s.SaveConfig(secretName, namespace, config.ConfigKey, secretData); err != nil {
		t.Error(err)
		return
	}

	if err := s.DeleteConfig(secretName, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := s.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Empty(t, list)
}

// TestCreateUpdateSecret test the creation and update of a secret
func TestCreateUpdateSecret(t *testing.T) {
	s, err := newSecret()
	if err != nil {
		t.Error(err)
		return
	}

	if err := s.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	secretName := name + "-" + s.Secrets.Files[0].Name

	secret, err := s.Get(secretName, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	expected := createSecret()
	assert.Equal(t, expected, secret)

	// update
	expected.StringData[config.ConfigKey] = secretUpdatedData
	if err := s.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	if _, err = s.Get(secretName, namespace); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, secretUpdatedData, expected.StringData[config.ConfigKey])
}

// TestDeleteSecrets test delete of a secret
func TestDeleteSecrets(t *testing.T) {
	s, err := newSecret()
	if err != nil {
		t.Error(err)
		return
	}

	if err := s.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := s.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list, 1)

	if err := s.Delete(name, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err = s.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list, 0)
}

func TestCreateUpdateSecretNotOwnedByKratos(t *testing.T) {
	s, err := newSecret()
	if err != nil {
		t.Error(err)
		return
	}

	if err := createNotKratosSecret(s); err != nil {
		t.Error(err)
		return
	}

	if err := s.CreateUpdate(name, namespace); assert.Error(t, err) {
		assert.Equal(t, "secret is not managed by kratos", err.Error())
	}
}

func TestDeleteSecretNotOwnedByKratos(t *testing.T) {
	s, err := newSecret()
	if err != nil {
		t.Error(err)
		return
	}

	if err := createNotKratosSecret(s); err != nil {
		t.Error(err)
		return
	}

	if err := s.Delete(name, namespace); assert.Error(t, err) {
		assert.Equal(t, "secret is not managed by kratos", err.Error())
	}
}
