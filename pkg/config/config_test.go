package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	goodConfig                 = "testdata/goodConfig.yml"
	badConfigResources         = "testdata/badConfigResources.yml"
	badConfigDeployment        = "testdata/badConfigDeployment.yml"
	badConfigDeploymentLabels  = "testdata/badConfigDeploymentLabels.yml"
	badConfigIngressLabels     = "testdata/badConfigIngressLabels.yml"
	badConfigCommonLabels      = "testdata/badConfigCommonLabels.yml"
	badConfigLabelsDuplication = "testdata/badConfigLabelsDuplication.yml"
	badConfigCommonAnnotations = "testdata/badConfigAnnotationsDuplication.yml"

	name                = "myapp"
	namespace           = "mynamespace"
	replicas      int32 = 1
	image               = "myimage"
	tag                 = "latest"
	port          int32 = 80
	ingresClass         = "nginx"
	clusterIssuer       = "letsencrypt"
	hostname            = "example.com"
)

func createConf() *Config {
	commonLabels := map[string]string{"environment": "dev"}
	commonAnnotations := map[string]string{"branch": "dev"}
	depLabels := map[string]string{"app": "myapp"}
	depAnnotations := map[string]string{"revision": "22"}
	ingLabels := map[string]string{"cloudflare": "enabled"}
	ingAnnotation := map[string]string{"hsts": "true"}
	return &Config{
		Common: Common{
			Labels:      commonLabels,
			Annotations: commonAnnotations,
		},
		Deployment: Deployment{
			Labels:      depLabels,
			Annotations: depAnnotations,
			Replicas:    replicas,
			Port:        port,
			Containers: []Container{
				{
					Name:  name,
					Image: image,
					Tag:   tag,
					Resources: Resources{
						Requests: ResourceType{
							CPU:    "25m",
							Memory: "32Mi",
						},
						Limits: ResourceType{
							CPU:    "50m",
							Memory: "64Mi",
						},
					},
				},
			},
		},
		Ingress: Ingress{
			Labels:        ingLabels,
			Annotations:   ingAnnotation,
			IngressClass:  ingresClass,
			ClusterIssuer: clusterIssuer,
			Hostnames:     []Hostnames{hostname},
		},
	}
}

func TestLoadConfig(t *testing.T) {
	config := &Config{}

	if err := config.Load(goodConfig); err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, createConf(), config)
}

func TestLoadConfigDeployment(t *testing.T) {
	config := &Config{}
	expected := "validation of configuration failed: Key: 'Config.Deployment.Port' Error:Field validation for 'Port' failed on the 'required' tag\nKey: 'Config.Deployment.Containers' Error:Field validation for 'Containers' failed on the 'required' tag"
	if err := config.Load(badConfigDeployment); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigDeploymentLabels(t *testing.T) {
	config := &Config{}
	expected := "validation of labels environment/branch failed: a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')"
	if err := config.Load(badConfigDeploymentLabels); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigCommonLabels(t *testing.T) {
	config := &Config{}
	expected := "validation of labels environment/test failed: a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')"
	if err := config.Load(badConfigCommonLabels); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigLabelsDuplication(t *testing.T) {
	config := &Config{}
	expected := "common labels \"environment\" cannot be duplicated in deployment labels"
	if err := config.Load(badConfigLabelsDuplication); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigAnnotationsDuplication(t *testing.T) {
	config := &Config{}
	expected := "common annotations \"branch\" cannot be duplicated in deployment annotations"
	if err := config.Load(badConfigCommonAnnotations); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigIngresslabels(t *testing.T) {
	config := &Config{}
	expected := "validation of labels cloudflare dns failed: a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')"
	if err := config.Load(badConfigIngressLabels); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigResources(t *testing.T) {
	config := &Config{}
	expected := "validation of configuration resources failed: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'\ncontainer: myapp -> requests cpu: 25f"
	if err := config.Load(badConfigResources); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}
