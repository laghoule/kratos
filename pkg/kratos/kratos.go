package kratos

import (
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
func New(confFile, kubeconfig string) (*Kratos, error) {
	conf := &config.Config{}

	if confFile != "" {
		if err := conf.Load(confFile); err != nil {
			return nil, err
		}
	}

	kClient, err := k8s.New(kubeconfig, conf)
	if err != nil {
		return nil, err
	}

	return &Kratos{
		Client: kClient,
		Config: conf,
	}, nil
}

// CheckDependency check if all dependency are met
func (k *Kratos) CheckDependency() error {
	// check if we meet k8s version requirement
	if err := k.CheckVersionDepency(); err != nil {
		return err
	}

	// dependency for deployment
	if k.Config.Deployment != nil {

		// validate clusterIssuer
		if cm, err := certmanager.New(*k.Client); err == nil {
			if err := cm.IsClusterIssuerExist(k.Client, k.Config.Deployment.Ingress.ClusterIssuer); err != nil {
				return err
			}
		} else {
			return err
		}

		// validate ingressClass
		if err := k.Client.IsIngressClassExist(k.Config.Deployment.Ingress.IngressClass); err != nil {
			return err
		}
	}

	return nil
}
