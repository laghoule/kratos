package k8s

import (
	//"testing"

	//"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	fakeName             = "fakeName"
	fakeNamespace        = "default"
	fakeResourceVersion  = "666"
	fakeDeployLabelName  = "kratos.io/name"
	fakeDeployLabelValue = fakeName
)

func newTest() *Client {
	client := &Client{}
	client.Clientset = fake.NewSimpleClientset()
	return client
}
