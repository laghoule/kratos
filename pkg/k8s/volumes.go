package k8s

import (
	"github.com/laghoule/kratos/pkg/config"

	corev1 "k8s.io/api/core/v1"
)

// getVolumesConfForContainer associate secrets and configmaps to volume in container, return list of VolumeMount & Volume
func getVolumesConfForContainer(name string, container *config.Container, conf *config.Config) ([]corev1.VolumeMount, []corev1.Volume) {
	volumesMount := []corev1.VolumeMount{}
	volumes := []corev1.Volume{}

	if conf.Secrets != nil {
		for _, file := range conf.Secrets.Files {
			if listContain(file.Mount.ExposedTo, container.Name) {
				volumesMount = append(volumesMount, corev1.VolumeMount{
					Name:      prefixSecretVolName + md5sum(file.Name),
					MountPath: file.Mount.Path,
					ReadOnly:  true,
				})
				volumes = append(volumes, corev1.Volume{
					Name: prefixSecretVolName + md5sum(file.Name),
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
