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
	labels := map[string]string{"mylabels": "myvalue"}
	annotations := map[string]string{"myannotations": "myvalue"}
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
		ConfigMaps: &config.ConfigMaps{
			Labels:      labels,
			Annotations: annotations,
			Files: []config.File{
				{
					Name: "configuration.yaml",
					Data: "my configuration data",
					Mount: config.Mount{
						Path: "/etc/config",
						ExposedTo: []string{
							name,
						},
					},
				},
			},
		},
		Secrets: &config.Secrets{
			Labels:      labels,
			Annotations: annotations,
			Files: []config.File{
				{
					Name: "secret.yaml",
					Data: "my secret data",
					Mount: config.Mount{
						Path: "/etc/secret",
						ExposedTo: []string{
							name,
						},
					},
				},
			},
		},
	}
}

// createSecretConfig return a secret object representing a release configuration
func createSecretConfig() *corev1.Secret {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return nil
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + config.ConfigSuffix,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					k8s.SecretLabelName: name + config.ConfigSuffix,
				}),
		},
		StringData: map[string]string{
			config.ConfigKey: "common:\n    labels:\n        app: myapp\n    annotations:\n        branch: dev\ndeployment:\n    replicas: 1\n    port: 80\n    containers:\n        - name: myapp\n          image: myimage\n          tag: latest\n          resources:\n            requests: {}\n            limits: {}\n          health:\n            live:\n                probe: /isLive\n                port: 80\n                initialDelaySeconds: 10\n                periodSeconds: 5\n            ready:\n                probe: /isReady\n                port: 80\n                initialDelaySeconds: 5\n                periodSeconds: 5\n    ingress:\n        ingressClass: nginx\n        clusterIssuer: letsencrypt\n        hostnames:\n            - example.com\nconfigmaps:\n    labels:\n        mylabels: myvalue\n    annotations:\n        myannotations: myvalue\n    files:\n        - name: configuration.yaml\n          data: my configuration data\n          mount:\n            path: /etc/config\n            exposedTo:\n                - myapp\nsecrets:\n    labels:\n        mylabels: myvalue\n    annotations:\n        myannotations: myvalue\n    files:\n        - name: secret.yaml\n          data: my secret data\n          mount:\n            path: /etc/secret\n            exposedTo:\n                - myapp\n",
		},
		Type: "Opaque",
	}
}

func TestSaveConfigFile(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.saveConfigToSecret(name+config.ConfigSuffix, namespace); err != nil {
		t.Error(err)
		return
	}

	result, err := k.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name+config.ConfigSuffix, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	expected := createSecretConfig()
	assert.Equal(t, expected, result)
}

func TestSaveConfigFileToDisk(t *testing.T) {
	k := new()
	if err := k.Config.Load(testdataInitFile); err != nil {
		t.Error(err)
		return
	}

	if err := k.saveConfigToSecret(name+config.ConfigSuffix, namespace); err != nil {
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
