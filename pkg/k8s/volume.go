package k8s

import (
	"github.com/laghoule/kratos/pkg/common"
	"github.com/laghoule/kratos/pkg/config"

	corev1 "k8s.io/api/core/v1"
)

// getVolumesConfForContainer associate secrets and configmaps to volume in container, return list of VolumeMount & Volume
func getVolumesConfForContainer(name string, container *config.Container, conf *config.Config) ([]corev1.VolumeMount, []corev1.Volume) {
	volumesMount := []corev1.VolumeMount{}
	volumes := []corev1.Volume{}

	if conf.ConfigMaps != nil {
		for _, file := range conf.ConfigMaps.Files {
			if common.ListContain(file.Mount.ExposedTo, container.Name) {
				volumesMount = append(volumesMount, corev1.VolumeMount{
					Name:      prefixConfigMapsVolName + common.MD5Sum16(file.Name),
					MountPath: file.Mount.Path,
					ReadOnly:  true,
				})
				volumes = append(volumes, corev1.Volume{
					Name: prefixConfigMapsVolName + common.MD5Sum16(file.Name),
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: name + "-" + file.Name,
							},
						},
					},
				})
			}
		}
	}

	if conf.Secrets != nil {
		for _, file := range conf.Secrets.Files {
			if common.ListContain(file.Mount.ExposedTo, container.Name) {
				volumesMount = append(volumesMount, corev1.VolumeMount{
					Name:      prefixSecretVolName + common.MD5Sum16(file.Name),
					MountPath: file.Mount.Path,
					ReadOnly:  true,
				})
				volumes = append(volumes, corev1.Volume{
					Name: prefixSecretVolName + common.MD5Sum16(file.Name),
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: name + "-" + file.Name,
						},
					},
				})
			}
		}
	}

	return volumesMount, volumes
}
