package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	goodConfig = "testdata/goodConfig.yml"
	badConfig  = "testdata/badConfig.yml"

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
					Name:      name,
					Image:     image,
					Tag:       tag,
					Resources: Resources{},
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

// TODO recheck badConfig.yml
func TestLoadBadConfig(t *testing.T) {
	config := &Config{}
	err := config.Load(badConfig)
	assert.Error(t, err)
}
