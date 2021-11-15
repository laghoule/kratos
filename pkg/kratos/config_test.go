package kratos

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createConf() *config.Config {
	return &config.Config{
		Deployment: &config.Deployment{
			Replicas: replicas,
			Containers: []config.Container{
				{
					Name:  name,
					Image: image,
					Tag:   tag,
					Port:  port,
				},
			},
		},
		Ingress: &config.Ingress{
			IngressClass:  ingresClass,
			ClusterIssuer: clusterIssuer,
			Hostnames:     []config.Hostnames{hostname},
		},
	}
}

func TestSaveConfigFile(t *testing.T) {
	c := new()
	c.Config = createConf()

	if err := c.saveConfigFileToSecret(name+kratosSuffixConfig, namespace); err != nil {
		t.Error(err)
		return
	}

	s, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name+kratosSuffixConfig, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	// TODO not enough for this test
	assert.Equal(t, name+kratosSuffixConfig, s.Name)
}

func TestSaveConfigFileToDisk(t *testing.T) {
	c := new()
	c.Config = createConf()

	if err := c.saveConfigFileToSecret(name+kratosSuffixConfig, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := c.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	if err := c.SaveConfigFileToDisk(name, namespace, os.TempDir()); err != nil {
		t.Error(err)
		return
	}

	result, err := os.ReadFile(filepath.Join(os.TempDir(), generatedInitFile))
	if err != nil {
		t.Error(err)
		return
	}

	expected, err := os.ReadFile(testdataInitFile)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, expected, result)
}

func TestCreateSecretString(t *testing.T) {
	s := createSecretString(name, namespace, configString)
	assert.Equal(t, configString, s.StringData[secretConfigKey])
}
