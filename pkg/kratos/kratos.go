package kratos

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/certmanager"
	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"
)

// Kratos contains info for deployment
type Kratos struct {
	*k8s.Client
	*config.Config
}

// New return a kratos struct
func New(conf, kubeconfig string) (*Kratos, error) {
	kclient, err := k8s.New(kubeconfig)
	if err != nil {
		return nil, err
	}

	// loading configuration
	confYAML := &config.Config{}
	if conf != "" {
		if err := confYAML.Load(conf); err != nil {
			return nil, err
		}
	}

	return &Kratos{
		Client: kclient,
		Config: confYAML,
	}, nil
}

// IsDependencyMeet check if all dependency are met
func (k *Kratos) IsDependencyMeet() error {
	var err error

	// check if we meet k8s version requirement
	if err := k.CheckVersionDepency(); err != nil {
		return err
	}

	// dependency for deployment
	if k.Config.Deployment != nil {

		cm := &certmanager.Certmanager{}
		if cm, err = certmanager.New(*k.Client); err != nil {
			return err
		}

		// validate clusterIssuer
		if err := cm.IsClusterIssuerExist(k.Client, k.Config.Deployment.Ingress.ClusterIssuer); err != nil {
			return err
		}

		// validate ingressClass
		if !k.Client.IsIngressClassExist(k.Config.Deployment.Ingress.IngressClass) {
			return fmt.Errorf("ingressClass %s not found", k.Config.Deployment.Ingress.ClusterIssuer)
		}
	}

	return nil
}
