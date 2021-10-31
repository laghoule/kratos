package k8s

import (
	"k8s.io/client-go/kubernetes/fake"
)


func newTest() *Client {
	client := &Client{}
	client.Clientset = fake.NewSimpleClientset()
	return client
}
