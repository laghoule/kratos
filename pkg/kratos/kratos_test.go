package kratos

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"
	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	goodConfig        = "../config/testdata/goodConfig.yml"
	badConfig         = "../config/testdata/badConfig.yml"
	generatedInitFile = "init.yaml"
	testdataInitFile  = "testdata/init.yaml"

	name                = "myapp"
	namespace           = "mynamespace"
	image               = "myimage"
	tag                 = "latest"
	replicas      int32 = 1
	port          int32 = 80
	ingresClass         = "nginx"
	clusterIssuer       = "letsencrypt"
	configString        = "my config"
	hostname            = "example.com"
)

func createConf() *config.Config {
	return &config.Config{
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
			Hostnames:     []config.Hostnames{hostname},
		},
	}
}

func new() *Kratos {
	kratos := &Kratos{
		Client: &k8s.Client{},
		Config: &config.Config{},
	}
	kratos.Clientset = fake.NewSimpleClientset()
	return kratos
}

func TestCreateInit(t *testing.T) {
	kratos := new()
	kratos.CreateInit(filepath.Join(os.TempDir(), generatedInitFile))

	expected, err := os.ReadFile(testdataInitFile)
	if err != nil {
		t.Error(err)
		return
	}

	result, err := os.ReadFile(filepath.Join(os.TempDir(), generatedInitFile))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, string(expected), string(result))
}

func TestSaveConfigFile(t *testing.T) {
	c := new()
	c.Config = createConf()

	if err := c.saveConfigFileToSecret(name+kratosSuffixConfig, namespace); err != nil {
		t.Error(err)
		return
	}

	s, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name+kratosSuffixConfig, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	// TODO not enough for this test
	assert.Equal(t, name+kratosSuffixConfig, s.Name)
}

func TestSreateSecretString(t *testing.T) {
	s := createSecretString(name, namespace, configString)
	assert.Equal(t, configString, s.StringData[secretConfigKey])
}

func TestSaveConfigFileToDisk(t *testing.T) {
	c := new()
	c.Config = createConf()

	if err := c.saveConfigFileToSecret(name+kratosSuffixConfig, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := c.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list.Items, 1)

	if err := c.SaveConfigFileToDisk(name, namespace, os.TempDir()); err != nil {
		t.Error(err)
		return
	}

	result, err := os.ReadFile(filepath.Join(os.TempDir(), generatedInitFile))
	if err != nil {
		t.Error(err)
		return
	}

	expected, err := os.ReadFile(testdataInitFile)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, expected, result)
}

func TestIsDependencyMeet(t *testing.T) {
	// TODO
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
