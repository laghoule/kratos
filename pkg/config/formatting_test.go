package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestFormatProbe(t *testing.T) {
	c := createDeploymentConf()

	lProbe := c.Deployment.Containers[0].FormatProbe(LiveProbe)
	rProbe := c.Deployment.Containers[0].FormatProbe(ReadyProbe)

	lExpected := &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: livePath,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: port,
				},
			},
		},
		InitialDelaySeconds: period,
		PeriodSeconds:       period,
	}

	rExpected := &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: readyPath,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: port,
				},
			},
		},
		InitialDelaySeconds: period,
		PeriodSeconds:       period,
	}

	assert.Equal(t, lExpected, lProbe)
	assert.Equal(t, rExpected, rProbe)

}

func TestFormatResources(t *testing.T) {
	c := createDeploymentConf()

	rs := c.Deployment.Containers[0].FormatResources()

	expected := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("25m"),
			corev1.ResourceMemory: resource.MustParse("32Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
	}

	assert.Equal(t, expected, rs)
}
