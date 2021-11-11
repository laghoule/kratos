package certmanager

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/k8s"

	certmanager "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
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

// CheckClusterIssuer check if a clusterIssuer exist
func (c *Certmanager)   CheckClusterIssuer(client *k8s.Client, clusterIssuer string) error {
	_, err := c.CertmanagerV1().ClusterIssuers().Get(context.Background(), clusterIssuer, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("clusterIsuser not found: %s", err)
		}
		return fmt.Errorf("getting clusterIssuer failed: %s", err)
	}

	return nil
}
