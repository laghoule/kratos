package k8s

import (
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	"golang.org/x/mod/semver"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
	depLabelName       = "kratos/deployment"
	cronLabelName      = "kratos/cronjob"
	requiredK8SVersion = "v1.19.0"
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

// formatResources format the resource from container configurations
func formatResources(container *config.Container) corev1.ResourceRequirements {
	if container == nil {
		return corev1.ResourceRequirements{
			Requests: corev1.ResourceList{},
			Limits:   corev1.ResourceList{},
		}
	}

	req := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{},
		Limits:   corev1.ResourceList{},
	}

	if container.Resources != nil {
		// requests
		if container.Resources.Requests != nil && container.Resources.Requests.CPU != "" {
			req.Requests[resCPU] = resource.MustParse(container.Resources.Requests.CPU)
		}
		if container.Resources.Requests != nil && container.Resources.Requests.Memory != "" {
			req.Requests[resMemory] = resource.MustParse(container.Resources.Requests.Memory)
		}

		// limits
		if container.Resources.Limits != nil && container.Resources.Limits.CPU != "" {
			req.Limits[resCPU] = resource.MustParse(container.Resources.Limits.CPU)
		}
		if container.Resources.Limits != nil && container.Resources.Limits.Memory != "" {
			req.Limits[resMemory] = resource.MustParse(container.Resources.Limits.Memory)
		}
	}

	return req
}
