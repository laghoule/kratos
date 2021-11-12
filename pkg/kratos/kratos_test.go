package kratos

import (
	"os"
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"
	"github.com/stretchr/testify/assert"

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
		Ingress: &config.Ingress{
			IngressClass:  ingresClass,
			ClusterIssuer: clusterIssuer,
			Hostnames:     hostnames,
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

func TestCreateInit(t *testing.T) {
	kratos := testNew()
	kratos.CreateInit("/tmp/init.yaml")

	expected, err := os.ReadFile("testdata/init.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	result, err := os.ReadFile("/tmp/init.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, expected, result)
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
