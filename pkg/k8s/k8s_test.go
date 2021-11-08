package k8s

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

const (
	kratosLabelName  = "app.kubernetes.io/managed-by"
	kratosLabelValue = "kratos"
	name             = "name"
	namespace        = "default"
	image            = "test"
	tagLatest        = "latest"
	tagV1            = "v1.0.0"
	containerHTTP    = 80
	constainerHTTPS  = 443
	clusterIssuer    = "letsencrypt"
	hostname         = "www.example.com"
)

func testNew() *Client {
	client := &Client{}
	client.Clientset = fake.NewSimpleClientset()
	return client
}

func TestCheckVersionDepency(t *testing.T) {
	client := testNew()
	if err := client.CheckVersionDepency(); err != nil {
		t.Error(err)
	}
}
