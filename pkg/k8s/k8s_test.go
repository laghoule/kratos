package k8s

import (
	"testing"

	"github.com/laghoule/kratos/pkg/common"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	kratosLabelName  = "app.kubernetes.io/managed-by"
	kratosLabelValue = "kratos"
	name             = "myapp"
	namespace        = "mynamespace"
	image            = "myimage"
	tagLatest        = "latest"
	tagV1            = "v1.0.0"
	containerHTTP    = 80
	clusterIssuer    = "letsencrypt"
	hostname         = "example.com"
	path             = "/"
	environment      = "dev"
	schedule         = "0 0 * * *"

	deploymentConfig = "../config/testdata/deploymentConfig.yml"
	cronjobConfig    = "../config/testdata/cronjobConfig.yml"
	secretConfig     = "../config/testdata/secretsConfig.yml"

	managedByLabel = "app.kubernetes.io/managed-by"
)

func new() *Client {
	c := &Client{}
	c.Clientset = fake.NewSimpleClientset()
	return c
}

func TestBoolPTR(t *testing.T) {
	expected := true
	assert.Equal(t, &expected, common.BoolPTR(true))
}

func TestCheckVersionDepency(t *testing.T) {
	// TODO: TestCheckVersionDepency
}

func TestMD5sum(t *testing.T) {
	expected := "74657374d41d8cd98f00b204e9800998ecf8427e"
	assert.Equal(t, expected, common.MD5Sum("test"))
}

func TestCheckKratosManaged(t *testing.T) {
	labels := map[string]string{managedByLabel: "kratos"}
	err := checkKratosManaged(labels)
	assert.NoError(t, err)
}

func TestCheckNotKratosManaged(t *testing.T) {
	labels := map[string]string{managedByLabel: "helm"}
	err := checkKratosManaged(labels)
	assert.Error(t, err)
}
