package kratos

import (
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"

	//"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	goodConfig = "../config/testdata/goodConfig.yml"
	badConfig  = "../config/testdata/badConfig.yml"

	name                = "myapp"
	namespace           = "mynamespace"
	image               = "myimage"
	tag                 = "latest"
	replicas      int32 = 1
	port          int32 = 80
	ingresClass         = "nginx"
	clusterIssuer       = "letsencrypt"
)

var (
	hostnames     = []config.Hostnames{"example.com"}
	configuration = &config.Config{
		Deployment: &config.Deployment{
			Replicas: replicas,
			Containers: []config.Container{
				{
					Name:  name,
					Image: image,
					Tag:   tag,
					Port:  port,
				},
			},
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

func TestCreate(t *testing.T) {
	// TODO
}

func TestList(t *testing.T) {
	// TODO
}

func TestDelete(t *testing.T) {
	// TODO
}
