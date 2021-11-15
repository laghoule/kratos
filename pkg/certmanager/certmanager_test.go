package certmanager

import (
	"testing"

	"github.com/laghoule/kratos/pkg/k8s"

	cmv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	fake "github.com/jetstack/cert-manager/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	cmSet := fake.NewSimpleClientset(createClusterIssuer())
	return &Certmanager{Interface: cmSet}
}

func TestCheckClusterIssuer(t *testing.T) {
	cmClient := new()
	found := false

	if cmClient.IsClusterIssuerExist(&k8s.Client{}, goodName) {
		found = true
	}

	assert.True(t, found)
}

func TestCheckBadClusterIssuer(t *testing.T) {
	cmClient := new()
	found := false

	if cmClient.IsClusterIssuerExist(&k8s.Client{}, badName) {
		found = true
	}

	assert.False(t, found)
}
