package kratos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"
	"github.com/stretchr/testify/assert"
	netv1 "k8s.io/api/networking/v1"

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

type fakeIngress struct{}

func (f fakeIngress) CheckIngressClassExist(name string) error       { return nil }
func (f fakeIngress) CreateUpdate(name, namespace string) error      { return nil }
func (f fakeIngress) Delete(name, namespace string) error            { return nil }
func (f fakeIngress) List(namespace string) ([]netv1.Ingress, error) { return []netv1.Ingress{}, nil }

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
			Ingress: fakeIngress{},
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

func TestCheckDependency(t *testing.T) {
	// TODO TestCheckDependency
}
