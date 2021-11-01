package k8s

import (
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
)

func newTest() *Client {
	client := &Client{}
	client.Clientset = fake.NewSimpleClientset()
	return client
}
