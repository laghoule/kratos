package kratos

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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
					Health: &config.Health{
						Live: &config.Check{
							Probe:               "/isLive",
							Port:                80,
							InitialDelaySeconds: 10,
							PeriodSeconds:       5,
						},
						Ready: &config.Check{
							Probe:               "/isReady",
							Port:                80,
							InitialDelaySeconds: 5,
							PeriodSeconds:       5,
						},
					},
				},
			},
			Ingress: &config.Ingress{
				IngressClass:  ingresClass,
				ClusterIssuer: clusterIssuer,
				Hostnames:     []string{hostname},
			},
		},
	}
}

// createSecret return a secret object
func createSecret() *corev1.Secret {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + configSuffix,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					"app":               name,
					k8s.SecretLabelName: name + configSuffix,
				}),
			Annotations: map[string]string{
				"branch": "dev",
			},
		},
		StringData: map[string]string{
			config.ConfigKey: "common:\n    labels:\n        app: myapp\n    annotations:\n        branch: dev\ndeployment:\n    replicas: 1\n    port: 80\n    containers:\n        - name: myapp\n          image: myimage\n          tag: latest\n          resources:\n            requests: {}\n            limits: {}\n          health:\n            live:\n                probe: /isLive\n                port: 80\n                initialDelaySeconds: 10\n                periodSeconds: 5\n            ready:\n                probe: /isReady\n                port: 80\n                initialDelaySeconds: 5\n                periodSeconds: 5\n    ingress:\n        ingressClass: nginx\n        clusterIssuer: letsencrypt\n        hostnames:\n            - example.com\n",
		},
		Type: "Opaque",
	}
}

func TestSaveConfigFile(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.saveConfigToSecret(name+configSuffix, namespace); err != nil {
		t.Error(err)
		return
	}

	result, err := k.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name+configSuffix, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	expected := createSecret()
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
