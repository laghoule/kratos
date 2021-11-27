package kratos

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createConf() *config.Config {
	return &config.Config{
		Common: config.Common{
			Labels: map[string]string{
				"app": "myapp",
			},
			Annotations: map[string]string{
				"branch": "dev",
			},
		},
		Deployment: config.Deployment{
			Replicas: replicas,
			Port:     port,
			Containers: []config.Container{
				{
					Name:  name,
					Image: image,
					Tag:   tag,
				},
			},
		},
		Ingress: config.Ingress{
			IngressClass:  ingresClass,
			ClusterIssuer: clusterIssuer,
			Hostnames:     []config.Hostnames{hostname},
		},
	}
}

func TestSaveConfigFile(t *testing.T) {
	c := new()
	c.Config = createConf()

	b, err := yaml.Marshal(c.Config)
	if err != nil {
		t.Error(err)
		return
	}

	expected := c.createSecretDataString(name+configSuffix, namespace, string(b))

	if err := c.saveConfigToSecret(name+configSuffix, namespace); err != nil {
		t.Error(err)
		return
	}

	result, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name+configSuffix, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, expected, result)
}

func TestSaveConfigFileToDisk(t *testing.T) {
	c := new()
	c.Config.Load(testdataInitFile)

	if err := c.saveConfigToSecret(name+configSuffix, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := c.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	if err := c.SaveConfigToDisk(name, namespace, os.TempDir()); err != nil {
		t.Error(err)
		return
	}

	result, err := os.ReadFile(filepath.Join(os.TempDir(), name+yamlExt))
	if err != nil {
		t.Error(err)
		return
	}

	expected, err := os.ReadFile(testdataInitFile)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, string(expected), string(result))
}

func TestCreateSecretString(t *testing.T) {
	c := new()
	s := c.createSecretDataString(name, namespace, configString)
	assert.Equal(t, configString, s.StringData[configKey])
}