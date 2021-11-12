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

var (
	clusterIssuer = &cmv1.ClusterIssuer{
		ObjectMeta: metav1.ObjectMeta{
			Name: goodName,
		},
	}
)

func testNew() *Certmanager {
	cmSet := fake.NewSimpleClientset(clusterIssuer)
	return &Certmanager{Interface: cmSet}
}

func TestCheckClusterIssuer(t *testing.T) {
	cmClient := testNew()

	err := cmClient.CheckClusterIssuer(&k8s.Client{}, goodName)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestCheckBadClusterIssuer(t *testing.T) {
	cmClient := testNew()
	err := cmClient.CheckClusterIssuer(&k8s.Client{}, badName)
	assert.Error(t, err)
}