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
	return &Config{
		Deployment: &Deployment{
			Replicas: replicas,
			Containers: []Container{
				{
					Name:  name,
					Image: image,
					Tag:   tag,
					Port:  port,
					Resources: Resources{},
				},
			},
		},
		Ingress: &Ingress{
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
