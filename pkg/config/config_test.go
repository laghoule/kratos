package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	goodConfig = "testdata/goodConfig.yml"
	badConfig  = "testdata/badConfig.yml"

	name                = "myApp"
	namespace           = "myNamespace"
	replicas      int32 = 1
	port          int32 = 80
	ingresClass         = "nginx"
	clusterIssuer       = "letsencrypt"
)

var (
	hostnames     = []Hostnames{"example.com", "www.example.com"}
	configuration = &Config{
		Name:      name,
		Namespace: namespace,
		Deployment: &Deployment{
			Replicas: replicas,
		},
		Service: &Service{
			Port: port,
		},
		Ingress: &Ingress{
			IngressClass:  ingresClass,
			ClusterIssuer: clusterIssuer,
			Hostnames:     hostnames,
			Port:          port,
		},
	}
)

func TestLoadConfig(t *testing.T) {
	config := &Config{}

	if err := config.Load(goodConfig); err != nil {
		t.Error(err)
	}

	assert.EqualValues(t, configuration, config)
}

func TestLoadBadConfig(t *testing.T) {
	config := &Config{}
	err := config.Load(badConfig)
	assert.Error(t, err)
}
