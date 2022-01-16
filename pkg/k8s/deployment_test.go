package k8s

import (
	"context"
	"testing"

	"github.com/laghoule/kratos/pkg/common"
	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
)

func newDeployment() (*Deployment, error) {
	conf := &config.Config{}

	if err := conf.Load(deploymentConfig); err != nil {
		return nil, err
	}

	return &Deployment{
		Clientset: fake.NewSimpleClientset(),
		Config:    conf,
	}, nil
}

func createDeployment() *appsv1.Deployment {
	var replicas int32 = 1
	depLabels := map[string]string{
		kratosLabelName: kratosLabelValue,
		DepLabelName:    name,
		"environment":   environment,
		"app":           name,
	}
	annotations := map[string]string{
		"branch":   environment,
		"revision": "22",
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      depLabels,
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					DepLabelName: name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        name,
					Namespace:   namespace,
					Labels:      depLabels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image + ":" + tagLatest,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: containerHTTP,
								},
							},
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
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/isLive",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 80,
										},
									},
								},
								InitialDelaySeconds: 3,
								PeriodSeconds:       3,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/isReady",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 80,
										},
									},
								},
								InitialDelaySeconds: 3,
								PeriodSeconds:       3,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "configmap-636f6e6669677572",
									MountPath: "/etc/config",
									ReadOnly:  true,
								},
								{
									Name:      "secret-7365637265742e79",
									MountPath: "/etc/secret",
									ReadOnly:  true,
								},
							},
						},
					},
					AutomountServiceAccountToken: common.BoolPTR(false),
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: common.BoolPTR(true),
					},
					Volumes: []corev1.Volume{
						{
							Name: "configmap-636f6e6669677572",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "myapp-configuration.yaml",
									},
								},
							},
						},
						{
							Name: "secret-7365637265742e79",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "myapp-secret.yaml",
								},
							},
						},
					},
				},
			},
		},
	}
}

// TestListNoDeployment test list of no deployment
func TestListNoDeployment(t *testing.T) {
	d, err := newDeployment()
	if err != nil {
		t.Error(err)
		return
	}

	list, err := d.List(namespace)
	if err != nil {
		return
	}

	assert.Empty(t, list)
}

// TestListDeployment test list of one deployment
func TestListDeployment(t *testing.T) {
	d, err := newDeployment()
	if err != nil {
		t.Error(err)
		return
	}

	if err := d.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := d.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, list)
}

// TestCreateDeploymentUpdateNotOwnedByKratos test update of deployment not owned by kratos
func TestCreateUpdateDeploymentNotOwnedByKratos(t *testing.T) {
	d, err := newDeployment()
	if err != nil {
		t.Error(err)
		return
	}

	dep := createDeployment()
	dep.Labels = nil

	_, err = d.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), dep, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if err := d.CreateUpdate(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "deployment is not managed by kratos")
	}
}

// TestCreateUpdateDeployment test creation of deployment
func TestCreateUpdateDeployment(t *testing.T) {
	d, err := newDeployment()
	if err != nil {
		t.Error(err)
		return
	}

	// create
	if err := d.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	dep, err := d.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, createDeployment(), dep)

	// update
	d.Containers[0].Tag = "v1.0.0"
	if err := d.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	dep, err = d.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, image+":"+tagV1, dep.Spec.Template.Spec.Containers[0].Image)
}

// TestDeleteDeploymentNotOwnedByKratos test delete of deployment not owned by kratos
func TestDeleteDeploymentNotOwnedByKratos(t *testing.T) {
	d, err := newDeployment()
	if err != nil {
		t.Error(err)
		return
	}

	dep := createDeployment()
	dep.Labels = nil

	_, err = d.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), dep, metav1.CreateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if err := d.Delete(name, namespace); assert.Error(t, err) {
		assert.Equal(t, err.Error(), "deployment is not managed by kratos")
	}
}

// TestDeleteDeployment test delete of deployment
func TestDeleteDeployment(t *testing.T) {
	d, err := newDeployment()
	if err != nil {
		t.Error(err)
		return
	}

	if err := d.CreateUpdate(name, namespace); err != nil {
		t.Error(err)
		return
	}

	dep, err := d.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEmpty(t, dep)

	if err := d.Delete(name, namespace); err != nil {
		t.Error(err)
		return
	}

	_, err = d.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		t.Error(err)
		return
	}

	assert.True(t, errors.IsNotFound(err))
}
