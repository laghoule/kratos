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
	if err := k.CheckVersionDepency(); err != nil {
		return err
	}

	if k.Config.Deployment != nil {
		if cm, err := certmanager.New(*k.Client); err == nil {
			if err := cm.CheckClusterIssuerExist(k.Client, k.Config.Deployment.Ingress.ClusterIssuer); err != nil {
				return err
			}
		} else {
			return err
		}

		if err := k.Client.CheckIngressClassExist(k.Config.Deployment.Ingress.IngressClass); err != nil {
			return err
		}
	}

	return nil
}

// loadConfigFromSecret load application config from secret
func (k *Kratos) loadConfigFromSecret(name, namespace string) error {
	secret, err := k.Get(name+config.ConfigSuffix, namespace)
	if err != nil {
		return fmt.Errorf("getting config from secret failed: %s", err)
	}

	var conf string
	if _, ok := secret.Data[config.ConfigKey]; ok {
		conf = string(secret.Data[config.ConfigKey])
	} else {
		if _, ok := secret.StringData[config.ConfigKey]; ok {
			conf = secret.StringData[config.ConfigKey]
		} else {
			return fmt.Errorf("getting config data from secret failed")
		}
	}

	if err := k.Config.LoadFromString(conf); err != nil {
		return err
	}

	return nil
}
