package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/config"
	"github.com/stretchr/testify/assert"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createCronjobs() *batchv1.CronJob {
	var retry int32 = 3
	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				kratosLabelName: kratosLabelValue,
				CronLabelName:   name,
				"environment":   environment,
				"type":          "long",
			},
			Annotations: map[string]string{
				"branch":   environment,
				"revision": "22",
			},
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          schedule,
			ConcurrencyPolicy: batchv1.ForbidConcurrent,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels: map[string]string{
						kratosLabelName: kratosLabelValue,
						CronLabelName:   name,
						"environment":   environment,
						"type":          "long",
					},
					Annotations: map[string]string{
						"branch":   environment,
						"revision": "22",
					},
				},
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name:      name,
							Namespace: namespace,
							Labels: map[string]string{
								kratosLabelName: kratosLabelValue,
								CronLabelName:   name,
								"environment":   environment,
								"type":          "long",
							},
							Annotations: map[string]string{
								"branch":   environment,
								"revision": "22",
							},
						},
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{
								{
									Name:  name,
									Image: image + ":" + tagLatest,
									Resources: corev1.ResourceRequirements{
										Requests: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("25m"),
											corev1.ResourceMemory: resource.MustParse("32Mi"),
										},
										Limits: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("50m"),
											corev1.ResourceMemory: resource.MustParse("64Mi"),
										},
									},
									VolumeMounts: []corev1.VolumeMount{},
								},
							},
							AutomountServiceAccountToken: boolPTR(false),
							SecurityContext: &corev1.PodSecurityContext{
								RunAsNonRoot: boolPTR(true),
							},
							Volumes: []corev1.Volume{},
						},
					},
					BackoffLimit: &retry,
				},
			},
		},
	}
}

// TestListCronjobs test the listing of cronjobs
func TestListCronjobs(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(cronjobConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateCronjob(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	list, err := c.ListCronjobs(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list, 1)
}

func TestCreateUpdateCronjobNotOwnedByKratos(t *testing.T) {
	c := new()
	conf := &config.Config{}

	cron := createCronjobs()
	cron.Labels = nil

	_, err := c.Clientset.BatchV1().CronJobs(namespace).Create(context.Background(), cron, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateCronjob(name, namespace, conf); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "cronjob is not owned by kratos")
	}
}

func TestCreateUpdateCronjob(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(cronjobConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateCronjob(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	result, err := c.Clientset.BatchV1().CronJobs(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	expected := createCronjobs()
	assert.Equal(t, expected, result)

	conf.Cronjob.Container.Tag = tagV1

	if err := c.CreateUpdateCronjob(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	result, err = c.Clientset.BatchV1().CronJobs(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, result.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image, image+":"+tagV1)
}

func TestDeleteCronjobNotOwnedByKratos(t *testing.T) {
	c := new()

	cron := createCronjobs()
	cron.Labels = nil

	_, err := c.Clientset.BatchV1().CronJobs(namespace).Create(context.Background(), cron, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if err := c.DeleteCronjob(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "cronjob is not owned by kratos")
	}
}

func TestDeleteCronjob(t *testing.T) {
	c := new()
	conf := &config.Config{}

	if err := conf.Load(cronjobConfig); err != nil {
		t.Error(err)
		return
	}

	if err := c.CreateUpdateCronjob(name, namespace, conf); err != nil {
		t.Error(err)
		return
	}

	result, err := c.Clientset.BatchV1().CronJobs(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	expected := createCronjobs()
	assert.Equal(t, expected, result)

	if err := c.DeleteCronjob(name, namespace); err != nil {
		t.Error(err)
		return
	}

	_, err = c.Clientset.BatchV1().CronJobs(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Error(err)
		return
	}

	assert.True(t, errors.IsNotFound(err))
}
