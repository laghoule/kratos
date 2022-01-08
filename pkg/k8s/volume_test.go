package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVolumesConfForContainer(t *testing.T) {
	d, err := newDeployment()
	if err != nil {
		t.Error(err)
		return
	}

	dep := createDeployment()

	volMounts, vols := getVolumesConfForContainer(name, &d.Deployment.Containers[0], d.Config)
	assert.Equal(t, dep.Spec.Template.Spec.Containers[0].VolumeMounts, volMounts)
	assert.Equal(t, dep.Spec.Template.Spec.Volumes, vols)
}
