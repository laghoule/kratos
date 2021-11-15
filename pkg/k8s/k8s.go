package k8s

import (
	"flag"
	"fmt"
	"path/filepath"

	"golang.org/x/mod/semver"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client is the Kubernetes client
type Client struct {
	Clientset  kubernetes.Interface
	RestConfig *rest.Config
}

const (
	appLabelName       = "app"
	requiredK8SVersion = "v1.19.0"
)

// New return a a Client
func New() (*Client, error) {
	var kubeconfig *string

	// TODO: FIX FLAG USE

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
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
