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

func new() *Kratos {
	kratos := &Kratos{
		Client: &k8s.Client{},
		Config: &config.Config{},
	}
	kratos.Clientset = fake.NewSimpleClientset()
	return kratos
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
