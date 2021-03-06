package certmanager

import (
	"testing"

	"github.com/laghoule/kratos/pkg/k8s"

	cmv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	fake "github.com/jetstack/cert-manager/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

const (
	goodName = "letsencrypt"
	badName  = "notletsencrypt"
)

func createClusterIssuer() *cmv1.ClusterIssuer {
	return &cmv1.ClusterIssuer{
		ObjectMeta: metav1.ObjectMeta{
			Name: goodName,
		},
	}
}

func new() *Certmanager {
	c := fake.NewSimpleClientset(createClusterIssuer())
	return &Certmanager{Interface: c}
}

func TestNew(t *testing.T) {
	k := k8s.Client{
		RestConfig: &rest.Config{},
	}

	_, err := New(k)
	assert.NoError(t, err)
}

func TestCheckClusterIssuer(t *testing.T) {
	c := new()

	err := c.CheckClusterIssuerExist(&k8s.Client{}, goodName)
	assert.NoError(t, err)
}

func TestCheckBadClusterIssuer(t *testing.T) {
	c := new()

	if err := c.CheckClusterIssuerExist(&k8s.Client{}, badName); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "clusterIssuer notletsencrypt not found")
	}

}
