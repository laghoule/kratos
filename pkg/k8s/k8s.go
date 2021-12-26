package k8s

import (
	"crypto/md5"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	"golang.org/x/mod/semver"

	"k8s.io/apimachinery/pkg/labels"
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
	// DepLabelName is label applied to deployment
	DepLabelName = "kratos/deployment"
	// CronLabelName is label applied to cronjob
	CronLabelName = "kratos/cronjob"
	// SecretLabelName is the label applied to secrets
	SecretLabelName = "kratos/secret"
	// ConfigmapsLabelName is the label applied to configmaps
	ConfigmapsLabelName = "kratos/configmaps"
	requiredK8SVersion  = "v1.19.0"
	prefixSecretVolName = "secret-"
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

// listContain return true if searchItem is found in the list of string
func listContain(list []string, searchItem string) bool {
	for _, item := range list {
		if item == searchItem {
			return true
		}
	}

	return false
}

// boolPTR return a bool pointer
func boolPTR(b bool) *bool {
	return &b
}

// md5sum return a md5sum from the input string
func md5sum(input string) string {
	hash := md5.New()
	return fmt.Sprintf("%x", hash.Sum([]byte(input)))
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
