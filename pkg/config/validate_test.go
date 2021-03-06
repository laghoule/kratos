package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mountPath = "/etc/cfg"
)

func TestLabelsValidation(t *testing.T) {
	label := map[string]string{"app": name}
	err := labelsValidation(label)
	assert.NoError(t, err)
}

func TestLabelsValidationError(t *testing.T) {
	label := map[string]string{"APP!": name}
	err := labelsValidation(label)
	assert.Error(t, err)
}

func TestMapKeyUniq(t *testing.T) {
	m1 := map[string]string{"label1": "value"}
	m2 := map[string]string{"label2": "value"}
	err := mapKeyUniq(m1, m2)
	assert.NoError(t, err)
}

func TestMapKeyUniqError(t *testing.T) {
	m1 := map[string]string{"label": "value"}
	m2 := map[string]string{"label": "value"}
	err := mapKeyUniq(m1, m2)
	assert.Error(t, err)
}

func TestConfigValidateConfig(t *testing.T) {
	c := createDeploymentConf()
	err := c.validateConfig()
	assert.NoError(t, err)
}

func TestConfigValidateDeploymentError(t *testing.T) {
	c := createDeploymentConf()
	c.Deployment.Labels["environment"] = environment
	err := c.validateConfig()
	assert.Error(t, err)
}

func TestConfigValidateCronjobError(t *testing.T) {
	c := createCronjobConf()
	c.Cronjob.Labels["environment"] = environment
	err := c.validateConfig()
	assert.Error(t, err)
}

func TestConfigValidateConfigMapsError(t *testing.T) {
	c := createDeploymentConf()
	c.ConfigMaps.Labels["environment"] = environment
	d := DataFiles(*c.ConfigMaps)
	err := d.validateConfig(c)
	assert.Error(t, err)
	err = c.validateConfig()
	assert.Error(t, err)
}

func TestConfigValidateSecretsError(t *testing.T) {
	c := createDeploymentConf()
	c.Secrets.Labels["environment"] = environment
	d := DataFiles(*c.Secrets)
	err := d.validateConfig(c)
	assert.Error(t, err)
	err = c.validateConfig()
	assert.Error(t, err)
}

func TestConfigValidateMountPathError(t *testing.T) {
	c := createDeploymentConf()
	c.ConfigMaps.Files[0].Mount.Path = "/tmp"
	c.Secrets.Files[0].Mount.Path = "/tmp"
	err := c.validateConfig()
	assert.Error(t, err)
}

func TestCommonValidateConfig(t *testing.T) {
	c := createCronjobConf()
	common := c.Common
	err := common.validateConfig()
	assert.NoError(t, err)
}

func TestCommonValidateConfigError(t *testing.T) {
	common := &Common{
		Labels: map[string]string{
			"APP!": "value",
		},
	}
	err := common.validateConfig()
	assert.Error(t, err)
}

func TestDeploymentValidateConfig(t *testing.T) {
	c := createDeploymentConf()
	common := c.Common
	dep := c.Deployment
	err := dep.validateConfig(common)
	assert.NoError(t, err)
}

func TestDeploymentValidateConfigError(t *testing.T) {
	c := createDeploymentConf()
	common := c.Common
	dep := c.Deployment
	dep.Labels = map[string]string{
		"APP!": "value",
	}
	err := dep.validateConfig(common)
	assert.Error(t, err)
}

func TestCronjobValidateConfig(t *testing.T) {
	c := createCronjobConf()
	common := c.Common
	cron := c.Cronjob
	err := cron.validateConfig(common)
	assert.NoError(t, err)
}

func TestCronjobValidateConfigError(t *testing.T) {
	c := createCronjobConf()
	common := c.Common
	cron := c.Cronjob
	cron.Schedule = "abc"
	err := cron.validateConfig(common)
	assert.Error(t, err)
}

func TestIngressValidateConfig(t *testing.T) {
	c := createDeploymentConf()
	common := c.Common
	ing := c.Ingress
	err := ing.validateConfig(common)
	assert.NoError(t, err)
}

func TestIngressValidateConfigError(t *testing.T) {
	c := createDeploymentConf()
	common := c.Common
	ing := c.Ingress
	ing.Labels = map[string]string{
		"APP!": "value",
	}
	err := ing.validateConfig(common)
	assert.Error(t, err)
}

func TestResourcesValidateConfig(t *testing.T) {
	c := createDeploymentConf()
	res := c.Deployment.Containers[0].Resources
	err := res.validateConfig(name)
	assert.NoError(t, err)
}

func TestResourcesValidateConfigError(t *testing.T) {
	c := createDeploymentConf()
	res := c.Deployment.Containers[0].Resources
	res.Limits = &ResourceType{
		CPU: "error",
	}
	err := res.validateConfig(name)
	assert.Error(t, err)
}

func TestResourceTypeValidateConfig(t *testing.T) {
	c := createDeploymentConf()
	resType := c.Deployment.Containers[0].Resources.Limits
	err := resType.validateConfig(name, "CPU")
	assert.NoError(t, err)
}

func TestResourceTypeValidateConfigERROR(t *testing.T) {
	c := createDeploymentConf()
	resType := c.Deployment.Containers[0].Resources.Limits
	resType.CPU = "error"
	err := resType.validateConfig(name, "CPU")
	assert.Error(t, err)
}

func TestSecretValidateConfig(t *testing.T) {
	c := createDeploymentConf()
	d := DataFiles(*c.Secrets)
	err := d.validateConfig(c)
	assert.NoError(t, err)
}

func TestSecretValidateConfigError(t *testing.T) {
	c := createDeploymentConf()
	secret := c.Secrets
	secret.Files[0].Mount.ExposedTo = []string{"error"}
	d := DataFiles(*c.Secrets)
	err := d.validateConfig(c)
	assert.Error(t, err)
}

func TestConfigMapsValidateConfig(t *testing.T) {
	c := createDeploymentConf()
	d := DataFiles(*c.ConfigMaps)
	err := d.validateConfig(c)
	assert.NoError(t, err)
}

func TestConfigMapsValidateConfigError(t *testing.T) {
	c := createDeploymentConf()
	cm := c.ConfigMaps
	cm.Files[0].Mount.ExposedTo = []string{"error"}
	d := DataFiles(*cm)
	err := d.validateConfig(c)
	assert.Error(t, err)
}

func TestValidateExposedTo(t *testing.T) {
	c := createDeploymentConf()
	err := validateExposedTo(c.Secrets.Files, c)
	assert.NoError(t, err)
}

func TestValidateExposedToError(t *testing.T) {
	c := createDeploymentConf()
	files := []File{
		{
			Mount: &Mount{
				Path: mountPath,
			},
		},
	}
	err := validateExposedTo(files, c)
	assert.Error(t, err)
}

func TestValidateMountPath(t *testing.T) {
	c := createDeploymentConf()
	err := c.validateMountPath()
	assert.NoError(t, err)
}

func TestValidateMountPathError(t *testing.T) {
	c := createDeploymentConf()
	c.Secrets.Files = []File{
		{
			Mount: &Mount{
				Path:      mountPath,
				ExposedTo: []string{name},
			},
		},
		{
			Mount: &Mount{
				Path:      mountPath,
				ExposedTo: []string{name},
			},
		},
	}
	err := c.validateMountPath()
	assert.Error(t, err)
}
