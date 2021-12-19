package k8s

import (
	"fmt"

	"golang.org/x/mod/semver"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client is the Kubernetes client
type Client struct {
	Clientset  kubernetes.Interface
	RestConfig *rest.Config
}

const (
	// TODO: Theses labels are not very useful, propably remove them in future
	DepLabelName        = "kratos/deployment"
	CronLabelName       = "kratos/cronjob"
	SecretLabelName     = "kratos/secret"
	ConfigmapsLabelName = "kratos/configmaps"
	requiredK8SVersion  = "v1.19.0"
)

// New return a a Client
func New(kubeconfig string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable get kubernetes client configuration: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Kubernetes client: %s", err)
	}

	return &Client{
		Clientset:  clientset,
		RestConfig: config,
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
