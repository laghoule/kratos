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

func createSecret() *corev1.Secret {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					appLabelName: name,
				}),
		},
		StringData: map[string]string{
			"mykey": "my secret data",
		},
		Type: "Opaque",
	}
}

func TestCreateSecret(t *testing.T) {
	c := new()
	s := createSecret()

	if err := c.CreateUpdateSecret(s, namespace); err != nil {
		t.Error(err)
		return
	}

	secret, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), s.Name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, s, secret)
}

func TestUpdateSecret(t *testing.T) {
	c := new()
	s := createSecret()

	if err := c.CreateUpdateSecret(s, namespace); err != nil {
		t.Error(err)
		return
	}

	secret, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), s.Name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	secret.StringData["mykey"] = "my updated secret data"

	if err := c.CreateUpdateSecret(secret, namespace); err != nil {
		t.Error(err)
		return
	}

	secret, err = c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), s.Name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "my updated secret data", secret.StringData["mykey"])
}

func TestDeleteSecret(t *testing.T) {
	c := new()
	s := createSecret()

	if err := c.CreateUpdateSecret(s, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := c.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	if err := c.DeleteSecret(s.Name, namespace); err != nil {
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

func TestGetSecret(t *testing.T) {
	c := new()
	s := createSecret()

	if err := c.CreateUpdateSecret(s, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := c.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	secret, err := c.GetSecret(name, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, s, secret)
}