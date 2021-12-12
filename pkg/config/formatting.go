package config

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// FormatProbe format the container healthcheck to return a *corev1.Probe
func (c *Container) FormatProbe(healthType string) *corev1.Probe {
	switch strings.ToLower(healthType) {

	// liveness
	case LiveProbe:
		if c.Health != nil && c.Health.Live != nil {
			probe := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: c.Health.Live.Probe,
						Port: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: c.Health.Live.Port,
						},
					},
				},
				InitialDelaySeconds: 1,
				PeriodSeconds:       1,
			}

			if c.Health.Live.InitialDelaySeconds != 0 {
				probe.InitialDelaySeconds = c.Health.Live.InitialDelaySeconds
			}

			if c.Health.Live.PeriodSeconds != 0 {
				probe.PeriodSeconds = c.Health.Live.PeriodSeconds
			}

			return probe
		}

	// readyness
	case ReadyProbe:
		if c.Health != nil && c.Health.Ready != nil {
			probe := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: c.Health.Ready.Probe,
						Port: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: c.Health.Ready.Port,
						},
					},
				},
				InitialDelaySeconds: 1,
				PeriodSeconds:       1,
			}

			if c.Health.Ready.InitialDelaySeconds != 0 {
				probe.InitialDelaySeconds = c.Health.Ready.InitialDelaySeconds
			}

			if c.Health.Ready.PeriodSeconds != 0 {
				probe.PeriodSeconds = c.Health.Ready.PeriodSeconds
			}

			return probe
		}

	default:
		return nil
	}

	return nil
}

// FormatResources format the resource of the container to return a corev1.ResourceRequirements
func (c *Container) FormatResources() corev1.ResourceRequirements {
	req := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{},
		Limits:   corev1.ResourceList{},
	}

	if c.Resources != nil {
		// requests
		if c.Resources.Requests != nil && c.Resources.Requests.CPU != "" {
			req.Requests["cpu"] = resource.MustParse(c.Resources.Requests.CPU)
		}
		if c.Resources.Requests != nil && c.Resources.Requests.Memory != "" {
			req.Requests["memory"] = resource.MustParse(c.Resources.Requests.Memory)
		}

		// limits
		if c.Resources.Limits != nil && c.Resources.Limits.CPU != "" {
			req.Limits["cpu"] = resource.MustParse(c.Resources.Limits.CPU)
		}
		if c.Resources.Limits != nil && c.Resources.Limits.Memory != "" {
			req.Limits["memory"] = resource.MustParse(c.Resources.Limits.Memory)
		}
	}

	return req
}
