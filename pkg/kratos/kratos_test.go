package kratos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"
	"github.com/stretchr/testify/assert"

	"k8s.io/client-go/kubernetes/fake"
)

const (
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
	hostname            = "example.com"
)

func new() *Kratos {
	conf := createConf()
	clientset := fake.NewSimpleClientset()
	return &Kratos{
		Client: &k8s.Client{
			Clientset: clientset,
			ConfigMaps: &k8s.ConfigMaps{
				Clientset: clientset,
				Config:    conf,
			},
			Cronjob: &k8s.Cronjob{
				Clientset: clientset,
				Config:    conf,
			},
			Deployment: &k8s.Deployment{
				Clientset: clientset,
				Config:    conf,
			},
			Ingress: &k8s.Ingress{
				Clientset: clientset,
				Config:    conf,
			},
			Secrets: &k8s.Secrets{
				Clientset: clientset,
				Config:    conf,
			},
			Service: &k8s.Service{
				Clientset: clientset,
				Config:    conf,
			},
		},
		Config: &config.Config{},
	}
}

func TestCreateInit(t *testing.T) {
	k := new()

	if err := k.CreateInit(filepath.Join(os.TempDir(), generatedInitFile)); err != nil {
		t.Error(err)
		return
	}

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

func TestIsDependencyMeet(t *testing.T) {
	// TODO TestIsDependencyMeet
}
