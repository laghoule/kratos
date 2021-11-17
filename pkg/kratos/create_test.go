package kratos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreate(t *testing.T) {
	c := new()
	c.Config = createConf()

	if err := c.Create(name, namespace); err != nil {
		t.Error(err)
		return
	}

	depList, err := c.Clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, depList.Items, 1)
	assert.Equal(t, name, depList.Items[0].Name)

	svcList, err := c.Clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, svcList.Items, 1)
	assert.Equal(t, name, svcList.Items[0].Name)

	ingList, err := c.Clientset.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, ingList.Items, 1)
	assert.Equal(t, name, ingList.Items[0].Name)

	secList, err := c.Clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, secList.Items, 1)
	assert.Equal(t, name+configSuffix, secList.Items[0].Name)
}
