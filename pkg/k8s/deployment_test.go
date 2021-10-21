package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	fakeDeployment = appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:              fakeName,
			Namespace:         fakeNamespace,
			ResourceVersion:   fakeResourceVersion,
			Generation:        1,
			CreationTimestamp: metav1.Time{},
			Labels: map[string]string{
				fakeDeployLabelName: fakeDeployLabelValue,
			},
		},
	}
)

func TestListNoDeployment(t *testing.T) {
	client := newTest()

	listDep, err := client.ListDeployment(fakeNamespace)
	if err != nil {
		return
	}

	assert.Empty(t, listDep)
}

func TestListDeployment(t *testing.T) {
	client := newTest()

	_, err := client.Clientset.AppsV1().Deployments(fakeNamespace).Create(context.Background(), &fakeDeployment, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
	}

	listDep, err := client.ListDeployment(fakeNamespace)
	if err != nil {
		t.Error(err)
	}

	assert.NotEmpty(t, listDep)
	assert.Len(t, listDep, 1)
}
