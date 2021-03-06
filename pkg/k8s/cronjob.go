package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/common"
	"github.com/laghoule/kratos/pkg/config"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Cronjob is the interface for cronjob
type Cronjob interface {
	CreateUpdate(string, string) error
	Delete(string, string) error
	List(string) ([]batchv1.CronJob, error)
}

// cronjob contain the kubernetes clientset and configuration of the release
type cronjob struct {
	Clientset kubernetes.Interface
	*config.Config
}

// checkOwnership check if it's safe to create, update or delete the cronjob
func (c *cronjob) checkOwnership(name, namespace string) error {
	cron, err := c.Clientset.BatchV1().CronJobs(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting cronjob failed: %s", err)
	}

	// managed by kratos
	if err := checkKratosManaged(cron.Labels); err == nil {
		if cron.Labels[CronLabelName] == name {
			return nil
		}
	}

	return fmt.Errorf("cronjob is not managed by kratos")
}

// CreateUpdate create or update a cronjobs
func (c *cronjob) CreateUpdate(name, namespace string) error {
	if err := c.checkOwnership(name, namespace); err != nil {
		return err
	}

	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.ManagedLabel)
	if err != nil {
		return fmt.Errorf("converting label failed: %s", err)
	}

	// merge labels
	if c.Common != nil && c.Common.Labels != nil {
		if err := mergeStringMaps(&c.Cronjob.Labels, c.Common.Labels, kratosLabel); err != nil {
			return fmt.Errorf("merging cronjob labels failed: %s", err)
		}
	}

	// merge annotations
	if c.Common != nil && c.Common.Annotations != nil {
		if err := mergeStringMaps(&c.Cronjob.Annotations, c.Common.Annotations); err != nil {
			return fmt.Errorf("merging cronjob annotations failed: %s", err)
		}
	}

	volumesMount, volumes := getVolumesConfForContainer(name, c.Cronjob.Container, c.Config)

	cronjobs := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				c.Cronjob.Labels,
				labels.Set{
					CronLabelName: name,
				},
			),
			Annotations: c.Cronjob.Annotations,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          c.Cronjob.Schedule,
			ConcurrencyPolicy: batchv1.ForbidConcurrent,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: labels.Merge(
						c.Cronjob.Labels,
						labels.Set{
							CronLabelName: name,
						},
					),
					Annotations: c.Cronjob.Annotations,
				},
				Spec: batchv1.JobSpec{
					BackoffLimit: &c.Cronjob.Retry,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name:      c.Cronjob.Container.Name,
							Namespace: namespace,
							Labels: labels.Merge(
								c.Cronjob.Labels,
								labels.Set{
									CronLabelName: name,
								},
							),
							Annotations: c.Cronjob.Annotations,
						},
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{
								{
									Name:         c.Cronjob.Container.Name,
									Image:        c.Cronjob.Container.Image + ":" + c.Cronjob.Container.Tag,
									Resources:    c.Cronjob.Container.FormatResources(),
									VolumeMounts: volumesMount,
								},
							},
							AutomountServiceAccountToken: common.PTR(false),
							SecurityContext: &corev1.PodSecurityContext{
								RunAsNonRoot: common.PTR(true),
							},
							Volumes: volumes,
						},
					},
				},
			},
		},
	}

	_, err = c.Clientset.BatchV1().CronJobs(namespace).Create(context.Background(), cronjobs, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err = c.Clientset.BatchV1().CronJobs(namespace).Update(context.Background(), cronjobs, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("updating cronjob failed: %s", err)
			}
		} else {
			return fmt.Errorf("creating cronjobs failes: %s", err)
		}
	}

	return nil
}

// Delete the specified cronjobs
func (c *cronjob) Delete(name, namespace string) error {
	if err := c.checkOwnership(name, namespace); err != nil {
		return err
	}

	if err := c.Clientset.BatchV1().CronJobs(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting cronjobs failed: %s", err)
	}

	return nil
}

// List cronjob of the specified namespace
func (c *cronjob) List(namespace string) ([]batchv1.CronJob, error) {
	list, err := c.Clientset.BatchV1().CronJobs(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: config.ManagedLabel,
	})
	if err != nil {
		return nil, fmt.Errorf("getting cronjobs list failed: %s", err)
	}

	return list.Items, nil
}
