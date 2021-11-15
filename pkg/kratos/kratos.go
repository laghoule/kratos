package kratos

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/certmanager"
	"github.com/laghoule/kratos/pkg/config"
	"github.com/laghoule/kratos/pkg/k8s"
)

// Kratosphere is the kratos interface
type Kratosphere interface {
	List(namespace string) error
	Create(name, namespace, image, tag, ingresClass, clusterIssuer string, hostnames []string, replicas, port int32) error
	Delete(name, namespace string)
}

// Kratos contains info for deployment
type Kratos struct {
	*k8s.Client
	*config.Config
}

// New return a kratos struct
func New(confFile string) (*Kratos, error) {
	kclient, err := k8s.New()
	if err != nil {
		return nil, err
	}

	// loading configuration
	confYAML := &config.Config{}
	if confFile != "" {
		if err := confYAML.Load(confFile); err != nil {
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
	// check if we meet k8s version requirement
	if err := k.CheckVersionDepency(); err != nil {
		return err
	}

	// validate clusterIssuer
	cm, err := certmanager.New(*k.Client)
	if err != nil {
		return err
	}

	if !cm.IsClusterIssuerExist(k.Client, k.Config.ClusterIssuer) {
		return fmt.Errorf("clusterIssuer %s not found", k.Config.ClusterIssuer)
	}

	// validate ingressClass
	if !k.Client.IsIngressClassExist(k.IngressClass) {
		return fmt.Errorf("ingressClass %s not found", k.ClusterIssuer)
	}

	return nil
}
