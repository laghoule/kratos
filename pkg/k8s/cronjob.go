package k8s

import (
	"context"
	"fmt"

	"github.com/laghoule/kratos/pkg/config"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/imdario/mergo"
)

// CreateUpdateCronjob create or update a cronjobs
func (c *Client) CreateUpdateCronjob(name, namespace string, conf *config.Config) error {
	kratosLabel, err := labels.ConvertSelectorToLabelsMap(config.DeployLabel)
	if err != nil {
		return fmt.Errorf("converting label failed: %s", err)
	}

	// merge common & cronjob labels
	if err := mergo.Map(&conf.Cronjob.Labels, conf.Common.Labels); err != nil {
		return fmt.Errorf("merging cronjob labels failed: %s", err)
	}

	// merge kratosLabels & cronjob labels
	if err := mergo.Map(&conf.Cronjob.Labels, map[string]string(kratosLabel)); err != nil {
		return fmt.Errorf("merging cronjob labels failed: %s", err)
	}

	// merge common & cronjob annotations
	if err := mergo.Map(&conf.Cronjob.Annotations, conf.Common.Annotations); err != nil {
		return fmt.Errorf("merging cronjob annotations failed: %s", err)
	}

	cronjobs := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: labels.Merge(
				conf.Cronjob.Labels,
				labels.Set{
					appLabelName: name,
				},
			),
			Annotations: conf.Cronjob.Annotations,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          conf.Cronjob.Schedule,
			ConcurrencyPolicy: batchv1.ForbidConcurrent,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: labels.Merge(
						conf.Cronjob.Labels,
						labels.Set{
							appLabelName: name,
						},
					),
					Annotations: conf.Cronjob.Annotations,
				},
				Spec: batchv1.JobSpec{
					BackoffLimit: &conf.Cronjob.Retry,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name:      conf.Cronjob.Container.Name,
							Namespace: namespace,
							Labels: labels.Merge(
								conf.Cronjob.Labels,
								labels.Set{
									appLabelName: name,
								},
							),
							Annotations: conf.Cronjob.Annotations,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:      conf.Cronjob.Container.Name,
									Image:     conf.Cronjob.Container.Image + ":" + conf.Cronjob.Container.Tag,
									Resources: formatResources(*conf.Cronjob.Container),
								},
							},
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

// DeleteCronjob delete the specified cronjobs
func (c *Client) DeleteCronjob(name, namespace string) error {
	if err := c.Clientset.BatchV1().CronJobs(namespace).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("deleting cronjobs failed: %s", err)
	}

	return nil
}
