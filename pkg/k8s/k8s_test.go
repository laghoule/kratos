package k8s

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

const (
	kratosLabelName         = "app.kubernetes.io/managed-by"
	kratosLabelValue        = "kratos"
	name                    = "myapp"
	namespace               = "mynamespace"
	image                   = "myimage"
	tagLatest               = "latest"
	tagV1                   = "v1.0.0"
	containerHTTP           = 80
	constainerHTTPS         = 443
	clusterIssuer           = "letsencrypt"
	hostname                = "example.com"
	path                    = "/"
)

func new() *Client {
	c := &Client{}
	c.Clientset = fake.NewSimpleClientset()
	return c
}

func TestCheckVersionDepency(t *testing.T) {
	c := new()
	if err := c.CheckVersionDepency(); err != nil {
		t.Error(err)
	}
}
