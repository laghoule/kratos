package certmanager

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/k8s"

	certmanager "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Certmanager contain the clientset
type Certmanager struct {
	certmanager.Interface
}

// New return a Cermanager object
func New(client k8s.Client) (*Certmanager, error) {
	cm, err := certmanager.NewForConfig(client.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("creating certmanager client failed: %s", err)
	}
	return &Certmanager{Interface: cm}, nil
}

// CheckClusterIssuerExist check if a clusterIssuer exist
func (c *Certmanager) CheckClusterIssuerExist(client *k8s.Client, clusterIssuer string) error {
	_, err := c.CertmanagerV1().ClusterIssuers().Get(context.Background(), clusterIssuer, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("clusterIssuer %s not found", clusterIssuer)
	}
	return nil
}
