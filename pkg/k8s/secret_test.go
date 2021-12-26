package k8s

import (
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	secretData        = "usename: patate\npassword: poil\n"
	secretUpdatedData = "my updated secret data"
	secretFileName    = "credentials.yaml"
)

// createSecret return a secret object
func createSecret() *corev1.Secret {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
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
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
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

// loadConfigCreateSecret load secretConfig file configuration and create secrets
func loadConfigCreateSecret(c *Client, conf *config.Config) error {
	if err := conf.Load(secretConfig); err != nil {
		return err
	}

	if err := c.CreateUpdateSecrets(name, namespace, conf); err != nil {
		return err
	}

	return nil
}

// loadSaveConfig load secretConfig and save release configuration in secret
func loadSaveConfig(c *Client, name, namespace string, conf *config.Config) error {
	if err := conf.Load(secretConfig); err != nil {
		return err
	}

	if err := c.SaveConfig(name, namespace, config.ConfigKey, secretData, conf); err != nil {
		return err
	}

	return nil
}

// TestSaveConfig test saving a release configuration
func TestSaveConfig(t *testing.T) {
	c := new()
	conf := &config.Config{}

	secretName := name + config.ConfigSuffix

	if err := loadSaveConfig(c, secretName, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	secret, err := c.GetSecret(secretName, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	expected := createConfigSecret()
	assert.Equal(t, expected, secret)
}

// TestDeleteConfig test deleting a release configuration
func TestDeleteConfig(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := loadSaveConfig(c, name+config.ConfigSuffix, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	if err := c.DeleteConfig(name+config.ConfigSuffix, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := c.listSecrets(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Empty(t, list.Items)
}

// TestCreateUpdateSecret test the creation and update of a secret
func TestCreateUpdateSecret(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := loadConfigCreateSecret(c, conf); err != nil {
		t.Error(err)
		return
	}

	secretName := name + "-" + conf.Secrets.Files[0].Name

	secret, err := c.GetSecret(secretName, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	expected := createSecret()
	assert.Equal(t, expected, secret)

	// update
	expected.StringData[config.ConfigKey] = secretUpdatedData
	if err := c.CreateUpdateSecrets(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	if _, err = c.GetSecret(secretName, namespace); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, secretUpdatedData, expected.StringData[config.ConfigKey])
}

// TestDeleteSecrets test delete of a secret
func TestDeleteSecrets(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := loadConfigCreateSecret(c, conf); err != nil {
		t.Error(err)
		return
	}

	list, err := c.listSecrets(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	if err := c.DeleteSecrets(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	list, err = c.listSecrets(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 0)
}

// TestGetSecret test getting a secret
func TestGetSecret(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := loadConfigCreateSecret(c, conf); err != nil {
		t.Error(err)
		return
	}

	secretName := name + "-" + conf.Secrets.Files[0].Name

	secret, err := c.GetSecret(secretName, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	expected := createSecret()

	assert.Equal(t, expected, secret)
}

func TestCreateUpdateSecretNotOwnedByKratos(t *testing.T) {
	// TODO: TestCreateUpdateSecretNotOwnedByKratos
}

func TestDeleteSecretNotOwnedByKratos(t *testing.T) {
	// TODO: TestDeleteSecretNotOwnedByKratos
}