package k8s

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	"golang.org/x/mod/semver"

	"github.com/imdario/mergo"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client contain the kubernetes clientset and the supported Kubernetes objects
type Client struct {
	Clientset  kubernetes.Interface
	RestConfig *rest.Config
	*ConfigMaps
	Cronjob
	Deployment
	Ingress
	*Secrets
	*Service
}

const (
	// DepLabelName is label applied to deployment
	DepLabelName = "kratos/deployment"
	// CronLabelName is label applied to cronjob
	CronLabelName = "kratos/cronjob"
	// SecretLabelName is the label applied to secrets
	SecretLabelName = "kratos/secret"
	// ConfigMapsLabelName is the label applied to configmaps
	ConfigMapsLabelName     = "kratos/configmaps"
	requiredK8SVersion      = "v1.19.0"
	prefixConfigMapsVolName = "configmap-"
	prefixSecretVolName     = "secret-"
)

// New return a a Client
func New(kubeconfig string, conf *config.Config) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable get kubernetes client configuration: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Kubernetes client: %s", err)
	}

	return &Client{
		Clientset:  clientset, // TODO: Check to eliminate the Clientset
		RestConfig: config,    // TODO: Needed for certmanager client, check to better integrate this
		ConfigMaps: &ConfigMaps{
			Clientset: clientset,
			Config:    conf,
		},
		Cronjob: &cronjob{
			Clientset: clientset,
			Config:    conf,
		},
		Deployment: &deployment{
			Clientset: clientset,
			Config:    conf,
		},
		Ingress: &ingress{
			Clientset: clientset,
			Config:    conf,
		},
		Secrets: &Secrets{
			Clientset: clientset,
			Config:    conf,
		},
		Service: &Service{
			Clientset: clientset,
			Config:    conf,
		},
	}, nil
}

// CheckVersionDepency check if depency are meet
func (c *Client) CheckVersionDepency() error {
	vers, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("getting Kubernetes version failed: %s", err)
	}

	if semver.Compare(vers.String(), requiredK8SVersion) < 0 {
		return fmt.Errorf("minimal Kubernetes version %s not meet", requiredK8SVersion)
	}

	return nil
}

// checkKratosManaged check if the object labels contains `app.kubernetes.io/managed-by=kratos` label
func checkKratosManaged(objLabels map[string]string) error {
	managedLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return fmt.Errorf("parsing labels failed: %s", err)
	}

	for cLabel, cValue := range objLabels {
		for kLabel, kValue := range managedLabel {
			if cLabel == kLabel {
				if cValue == kValue {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("not managed by kratos")
}

// mergeStringMaps to destination string map
func mergeStringMaps(dst *map[string]string, sources ...map[string]string) error {
	for _, src := range sources {
		if err := mergo.Map(dst, src); err != nil {
			return err
		}
	}

	return nil
}
