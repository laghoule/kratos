package kratos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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

type fakeConfigmaps struct{}
type fakeCronjob struct{}
type fakeDeployment struct{}
type fakeIngress struct{}

func (f fakeConfigmaps) CreateUpdate(name, namespace string) error { return nil }
func (f fakeConfigmaps) Delete(name, namespace string) error       { return nil }
func (f fakeConfigmaps) List(namespace string) ([]corev1.ConfigMap, error) {
	return []corev1.ConfigMap{}, nil
}

func (f fakeCronjob) CreateUpdate(name, namespace string) error { return nil }
func (f fakeCronjob) Delete(name, namespace string) error       { return nil }
func (f fakeCronjob) List(namespace string) ([]batchv1.CronJob, error) {
	return []batchv1.CronJob{}, nil
}

func (f fakeDeployment) CreateUpdate(name, namespace string) error { return nil }
func (f fakeDeployment) Delete(name, namespace string) error       { return nil }
func (f fakeDeployment) List(namespace string) ([]appsv1.Deployment, error) {
	return []appsv1.Deployment{}, nil
}

func (f fakeIngress) CheckIngressClassExist(name string) error  { return nil }
func (f fakeIngress) CreateUpdate(name, namespace string) error { return nil }
func (f fakeIngress) Delete(name, namespace string) error       { return nil }
func (f fakeIngress) List(namespace string) ([]netv1.Ingress, error) {
	return []netv1.Ingress{}, nil
}

func new() *Kratos {
	conf := createConf()
	clientset := fake.NewSimpleClientset()
	return &Kratos{
		Client: &k8s.Client{
			Clientset: clientset,
			ConfigMaps: fakeConfigmaps{},
			Cronjob:    fakeCronjob{},
			Deployment: fakeDeployment{},
			Ingress:    fakeIngress{},
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
