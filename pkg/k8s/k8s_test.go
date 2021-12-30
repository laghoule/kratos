package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	configMapsConfig = "../config/testdata/configmapsConfig.yml"

	managedByLabel = "app.kubernetes.io/managed-by"
)

func TestCheckVersionDepency(t *testing.T) {
	// TODO: TestCheckVersionDepency
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
