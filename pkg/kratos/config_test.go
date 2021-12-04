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
		Common: &config.Common{
			Labels: map[string]string{
				"app": name,
			},
			Annotations: map[string]string{
				"branch": "dev",
			},
		},
		Deployment: &config.Deployment{
			Replicas: replicas,
			Port:     port,
			Containers: []config.Container{
				{
					Name:  name,
					Image: image,
					Tag:   tag,
					Resources: &config.Resources{
						Requests: &config.ResourceType{},
						Limits:   &config.ResourceType{},
					},
				},
			},
		},
		Ingress: &config.Ingress{
			IngressClass:  ingresClass,
			ClusterIssuer: clusterIssuer,
			Hostnames:     []string{hostname},
		},
	}
}

func TestSaveConfigFile(t *testing.T) {
	k := new()
	k.Config = createConf()

	b, err := yaml.Marshal(k.Config)
	if err != nil {
		t.Error(err)
		return
	}

	expected := k.createSecretDataString(name+configSuffix, namespace, string(b))

	if err := k.saveConfigToSecret(name+configSuffix, namespace); err != nil {
		t.Error(err)
		return
	}

	result, err := k.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name+configSuffix, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, expected, result)
}

func TestSaveConfigFileToDisk(t *testing.T) {
	k := new()
	k.Config.Load(testdataInitFile)

	if err := k.saveConfigToSecret(name+configSuffix, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := k.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	if err := k.SaveConfigToDisk(name, namespace, os.TempDir()); err != nil {
		t.Error(err)
		return
	}

	result, err := os.ReadFile(filepath.Join(os.TempDir(), name+YamlExt))
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
	k := new()
	s := k.createSecretDataString(name, namespace, configString)
	assert.Equal(t, configString, s.StringData[configKey])
}
