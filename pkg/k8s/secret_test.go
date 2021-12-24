package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// createSecret return a secret object
func createSecret() *corev1.Secret {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	name := "credentials.yaml"

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					"app":           "myapp",
					"environment":   environment,
					SecretLabelName: name,
				}),
			Annotations: map[string]string{
				"branch": environment,
			},
		},
		StringData: map[string]string{
			name: "usename: patate\npassword: poil\n",
		},
		Type: "Opaque",
	}
}

func createK8SSecret(c *Client, conf *config.Config) error {
	if err := conf.Load(secretConfig); err != nil {
		return err
	}

	if err := c.CreateUpdateSecrets(namespace, conf); err != nil {
		return err
	}

	return nil
}

func TestSaveConfig(t *testing.T) {
	// TODO: TestSaveConfig
}

func TestDeleteConfig(t *testing.T) {
	// TODO: TestDeleteConfig
}

// TestCreateUpdateSecret test the creation and update of a secret
func TestCreateUpdateSecret(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := createK8SSecret(c, conf); err != nil {
		t.Error(err)
		return
	}

	name := conf.Secrets.Files[0].Name

	secret, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	expected := createSecret()
	assert.Equal(t, expected, secret)

	// update
	expected.StringData[config.ConfigKey] = "my updated secret data"
	if err := c.CreateUpdateSecrets(namespace, conf); err != nil {
		t.Error(err)
		return
	}

	if _, err = c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "my updated secret data", expected.StringData[config.ConfigKey])
}

func TestDeleteSecrets(t *testing.T) {
	// TODO: TestDeleteSecrets
}

// TestDeleteSecret test delete of a secret
func TestDeleteSecret(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := createK8SSecret(c, conf); err != nil {
		t.Error(err)
		return
	}

	list, err := c.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	if err := c.DeleteSecrets(namespace, conf); err != nil {
		t.Error(err)
		return
	}

	list, err = c.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
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

	if err := createK8SSecret(c, conf); err != nil {
		t.Error(err)
		return
	}

	name := conf.Secrets.Files[0].Name

	secret, err := c.GetSecret(name, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	expected := createSecret()

	assert.Equal(t, expected, secret)
}
