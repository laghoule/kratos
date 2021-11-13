package kratos

import (
	"context"
	"os"
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"
	"github.com/stretchr/testify/assert"

	//"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	assert.Equal(t, string(expected), string(result))
}

func TestSaveConfigFile(t *testing.T) {
	c := testNew()
	c.Config = configuration

	if err := c.saveConfigFile(name+kratosSuffixConfig, namespace); err != nil {
		t.Error(err)
		return
	}

	s, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name+kratosSuffixConfig, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	// TODO not enough for test
	assert.Equal(t, name+kratosSuffixConfig, s.Name)
}

func TestSreateSecretString(t *testing.T) {
	s := createSecretString(name, namespace, "my config")
	assert.Equal(t, "my config", s.StringData[kratosConfigKey])
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
