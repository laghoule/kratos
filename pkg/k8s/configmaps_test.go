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
	configData        = "setting1: patate\nsetting2: poil\n"
	configUpdatedData = "my updated config data"
	configFileName    = "settings.yaml"
)

func newConfigMaps() (*ConfigMaps, error) {
	conf := &config.Config{}

	if err := conf.Load(configMapsConfig); err != nil {
		return nil, err
	}

	return &ConfigMaps{
		Clientset: fake.NewSimpleClientset(),
		Config:    conf,
	}, nil
}

// createConfigmap return a configmap object
func createConfigMap() *corev1.ConfigMap {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return nil
	}

	cmName := name + "-" + configFileName

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: namespace,
			Labels: labels.Merge(
				kratosLabel,
				labels.Set{
					"app":               name,
					"environment":       environment,
					ConfigMapsLabelName: cmName,
				}),
			Annotations: map[string]string{
				"branch": environment,
			},
		},
		Data: map[string]string{
			configFileName: configData,
		},
	}
}

func createNotKratosConfigMaps(c *ConfigMaps) error {
	configmap := createConfigMap()
	configmap.Labels = nil

	if _, err := c.Clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), configmap, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

// TestCreateUpdateConfigMaps test the creation and update of a configmaps
func TestCreateUpdateConfigMaps(t *testing.T) {
	c, err := newConfigMaps()
	if err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	cmName := name + "-" + c.ConfigMaps.Files[0].Name

	configmap, err := c.get(cmName, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	expected := createConfigMap()
	assert.Equal(t, expected, configmap)

	// update
	expected.Data[configFileName] = configUpdatedData
	if err := c.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	if _, err = c.get(cmName, namespace); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, configUpdatedData, expected.Data[configFileName])
}

// TestDeleteConfigMaps test delete of a configmaps
func TestDeleteConfigMaps(t *testing.T) {
	c, err := newConfigMaps()
	if err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	if err := c.Delete(name, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := c.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list, 0)
}

func TestCreateUpdateConfigMapsNotOwnedByKratos(t *testing.T) {
	c, err := newConfigMaps()
	if err != nil {
		t.Error(err)
		return
	}

	if err := createNotKratosConfigMaps(c); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdate(name, namespace); assert.Error(t, err) {
		assert.Equal(t, "configmap is not managed by kratos", err.Error())
	}
}

func TestDeleteConfigMapsNotOwnedByKratos(t *testing.T) {
	c, err := newConfigMaps()
	if err != nil {
		t.Error(err)
		return
	}

	if err := createNotKratosConfigMaps(c); err != nil {
		t.Error(err)
		return
	}

	if err := c.Delete(name, namespace); assert.Error(t, err) {
		assert.Equal(t, "configmap is not managed by kratos", err.Error())
	}
}
