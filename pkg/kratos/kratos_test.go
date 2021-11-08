package kratos

import (
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	goodConfig = "../config/testdata/goodConfig.yml"
	badConfig  = "../config/testdata/badConfig.yml"

	name                = "myApp"
	namespace           = "myNamespace"
	replicas      int32 = 1
	port          int32 = 80
	ingresClass         = "nginx"
	clusterIssuer       = "letsencrypt"
)

var (
	hostnames     = []config.Hostnames{"example.com", "www.example.com"}
	configuration = &config.Config{
		Name:      name,
		Namespace: namespace,
		Deployment: &config.Deployment{
			Replicas: replicas,
		},
		Service: &config.Service{
			Port: port,
		},
		Ingress: &config.Ingress{
			IngressClass:  ingresClass,
			ClusterIssuer: clusterIssuer,
			Hostnames:     hostnames,
			Port:          port,
		},
	}
)

func testNew() *Kratos {
	kratos := &Kratos{
		Client: &k8s.Client{},
		Config: &config.Config{},
	}
	kratos.Clientset = fake.NewSimpleClientset()
	return kratos
}

func TestUseGoodConfig(t *testing.T) {
	client := testNew()

	if err := client.UseConfig(goodConfig); err != nil {
		t.Error(err)
	}

	assert.Equal(t, configuration, client.Config)
}

func TestUseBadConfig(t *testing.T) {
	client := testNew()
	err := client.UseConfig(badConfig)
	assert.Error(t, err)
}
