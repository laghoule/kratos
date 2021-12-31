package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	cronjobConfig              = "testdata/cronjobConfig.yml"
	cronjobConfigMinimal       = "testdata/cronjobConfigMinimal.yml"
	deploymentConfig           = "testdata/deploymentConfig.yml"
	deploymentConfigMinimal    = "testdata/deploymentConfigMinimal.yml"
	badConfigResources         = "testdata/badConfigResources.yml"
	badConfigDeployment        = "testdata/badConfigDeployment.yml"
	badConfigDeploymentLabels  = "testdata/badConfigDeploymentLabels.yml"
	badConfigIngressLabels     = "testdata/badConfigIngressLabels.yml"
	badConfigCommonLabels      = "testdata/badConfigCommonLabels.yml"
	badConfigLabelsDuplication = "testdata/badConfigLabelsDuplication.yml"
	badConfigCommonAnnotations = "testdata/badConfigAnnotationsDuplication.yml"

	name                = "myapp"
	replicas      int32 = 1
	image               = "myimage"
	tag                 = "latest"
	port          int32 = 80
	ingresClass         = "nginx"
	clusterIssuer       = "letsencrypt"
	hostname            = "example.com"
	schedule            = "0 0 * * *"
	retry         int32 = 3
	environment         = "dev"
	period        int32 = 3
	livePath            = "/isLive"
	readyPath           = "/isReady"
)

func createDeploymentConf() *Config {
	commonLabels := map[string]string{"environment": environment}
	commonAnnotations := map[string]string{"branch": environment}
	depLabels := map[string]string{"app": name}
	depAnnotations := map[string]string{"revision": "22"}
	ingLabels := map[string]string{"cloudflare": "enabled"}
	ingAnnotation := map[string]string{"hsts": "true"}
	labels := map[string]string{"mylabels": "myvalue"}
	annotations := map[string]string{"myannotations": "myvalue"}
	return &Config{
		Common: &Common{
			Labels:      commonLabels,
			Annotations: commonAnnotations,
		},
		Deployment: &Deployment{
			Labels:      depLabels,
			Annotations: depAnnotations,
			Replicas:    replicas,
			Port:        port,
			Containers: []Container{
				{
					Name:  name,
					Image: image,
					Tag:   tag,
					Resources: &Resources{
						Requests: &ResourceType{
							CPU:    "25m",
							Memory: "32Mi",
						},
						Limits: &ResourceType{
							CPU:    "50m",
							Memory: "64Mi",
						},
					},
					Health: &Health{
						Live: &Check{
							Probe:               livePath,
							Port:                port,
							InitialDelaySeconds: period,
							PeriodSeconds:       period,
						},
						Ready: &Check{
							Probe:               readyPath,
							Port:                port,
							InitialDelaySeconds: period,
							PeriodSeconds:       period,
						},
					},
				},
			},
			Ingress: &Ingress{
				Labels:        ingLabels,
				Annotations:   ingAnnotation,
				IngressClass:  ingresClass,
				ClusterIssuer: clusterIssuer,
				Hostnames:     []string{hostname},
			},
		},
		ConfigMaps: &ConfigMaps{
			Labels:      labels,
			Annotations: annotations,
			Files: []File{
				{
					Name: "configuration.yaml",
					Data: "my configuration data",
					Mount: Mount{
						Path: "/etc/config",
						ExposedTo: []string{
							name,
						},
					},
				},
			},
		},
		Secrets: &Secrets{
			Labels:      labels,
			Annotations: annotations,
			Files: []File{
				{
					Name: "secret.yaml",
					Data: "my secret data",
					Mount: Mount{
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

func createDeploymentConfMinimal() *Config {
	return &Config{
		Common: &Common{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Deployment: &Deployment{
			Replicas: replicas,
			Port:     port,
			Containers: []Container{
				{
					Name:  name,
					Image: image,
					Tag:   tag,
				},
			},
			Ingress: &Ingress{
				IngressClass:  ingresClass,
				ClusterIssuer: clusterIssuer,
				Hostnames:     []string{hostname},
			},
		},
	}
}

func createCronjobConf() *Config {
	commonLabels := map[string]string{"environment": environment}
	commonAnnotations := map[string]string{"branch": environment}
	return &Config{
		Common: &Common{
			Labels:      commonLabels,
			Annotations: commonAnnotations,
		},
		Cronjob: &Cronjob{
			Labels: map[string]string{
				"type": "long",
			},
			Annotations: map[string]string{
				"revision": "22",
			},
			Schedule: schedule,
			Retry:    retry,
			Container: &Container{
				Name:  name,
				Image: image,
				Tag:   tag,
				Resources: &Resources{
					Requests: &ResourceType{
						CPU:    "25m",
						Memory: "32Mi",
					},
					Limits: &ResourceType{
						CPU:    "50m",
						Memory: "64Mi",
					},
				},
			},
		},
	}
}

func createCronjobConfMinimal() *Config {
	return &Config{
		Common: &Common{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Cronjob: &Cronjob{
			Schedule: schedule,
			Retry:    retry,
			Container: &Container{
				Name:  name,
				Image: image,
				Tag:   tag,
			},
		},
	}
}

func TestLoadConfigDeployment(t *testing.T) {
	c := &Config{}

	if err := c.Load(deploymentConfig); err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, createDeploymentConf(), c)
}

func TestLoadConfigDeploymentMinimal(t *testing.T) {
	c := &Config{}

	if err := c.Load(deploymentConfigMinimal); err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, createDeploymentConfMinimal(), c)
}

func TestLoadConfigCronjob(t *testing.T) {
	c := &Config{}
	if err := c.Load(cronjobConfig); err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, createCronjobConf(), c)
}

func TestLoadConfigCronjobMinimal(t *testing.T) {
	c := &Config{}
	if err := c.Load(cronjobConfigMinimal); err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, createCronjobConfMinimal(), c)
}

func TestLoadConfigBadDeployment(t *testing.T) {
	c := &Config{}
	expected := "validation of configuration failed: Key: 'Config.Deployment.Port' Error:Field validation for 'Port' failed on the 'required' tag\nKey: 'Config.Deployment.Containers' Error:Field validation for 'Containers' failed on the 'required' tag"
	if err := c.Load(badConfigDeployment); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigBadDeploymentLabels(t *testing.T) {
	c := &Config{}
	expected := "validation of labels environment/branch failed: a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')"
	if err := c.Load(badConfigDeploymentLabels); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigBadCommonLabels(t *testing.T) {
	c := &Config{}
	expected := "validation of labels environment/test failed: a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')"
	if err := c.Load(badConfigCommonLabels); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigLabelsDuplication(t *testing.T) {
	c := &Config{}
	expected := "common labels/annotations \"environment\" cannot be duplicated in others sections"
	if err := c.Load(badConfigLabelsDuplication); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigAnnotationsDuplication(t *testing.T) {
	c := &Config{}
	expected := "common labels/annotations \"branch\" cannot be duplicated in others sections"
	if err := c.Load(badConfigCommonAnnotations); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigBadIngresslabels(t *testing.T) {
	c := &Config{}
	expected := "validation of labels cloudflare dns failed: a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')"
	if err := c.Load(badConfigIngressLabels); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestLoadConfigBadResources(t *testing.T) {
	c := &Config{}
	expected := "validation of configuration resources failed: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'\ncontainer: myapp -> requests cpu: 25f"
	if err := c.Load(badConfigResources); assert.Error(t, err) {
		assert.Equal(t, expected, err.Error())
	}
}

func TestListContainerNames(t *testing.T) {
	// TODO: TestListContainerNames
}
