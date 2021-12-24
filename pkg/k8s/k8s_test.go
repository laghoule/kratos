package k8s

import (
	"testing"

	//"github.com/stretchr/testify/assert"
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
)

func new() *Client {
	c := &Client{}
	c.Clientset = fake.NewSimpleClientset()
	return c
}

func TestCheckVersionDepency(t *testing.T) {
	// TODO: TestCheckVersionDepency
}
