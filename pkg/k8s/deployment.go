package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/kratos"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListDeployment list deployment of k8dep labels
func (c *Client) ListDeployment(ns string) ([]appsv1.Deployment, error) {
	list, err := c.Clientset.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{
		LabelSelector: kratos.DeployLabel,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting list of deployment: %s", err)
	}

	return list.Items, nil
}