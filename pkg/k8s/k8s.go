package k8s

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	"golang.org/x/mod/semver"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client contain the kubernetes clientset and the supported Kubernetes objects
type Client struct {
	Clientset  kubernetes.Interface
	RestConfig *rest.Config
	*Cronjob
	*Deployment
	*Ingress
	*Secret
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
	ConfigMapsLabelName = "kratos/configmaps"
	requiredK8SVersion  = "v1.19.0"
	prefixSecretVolName = "secret-"
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
		Cronjob: &Cronjob{
			Clientset: clientset,
			Config:    conf,
		},
		Deployment: &Deployment{
			Clientset: clientset,
			Config:    conf,
		},
		Ingress: &Ingress{
			Clientset: clientset,
			Config:    conf,
		},
		Secret: &Secret{
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
	managedLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
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
