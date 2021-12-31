package kratos

import (
	"testing"

	"github.com/laghoule/kratos/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.Create(name, namespace); err != nil {
		t.Error(err)
		return
	}

	// deployment
	depList, err := k.Client.Deployment.List(namespace)
	assert.NoError(t, err)
	assert.Equal(t, name, depList[0].Name)

	// service
	svcList, err := k.Client.Service.List(namespace)
	assert.NoError(t, err)
	assert.Equal(t, name, svcList[0].Name)

	// ingress
	ingList, err := k.Client.Ingress.List(namespace)
	assert.NoError(t, err)
	assert.Equal(t, name, ingList[0].Name)

	// configmaps
	cmList, err := k.Client.ConfigMaps.List(namespace)
	assert.NoError(t, err)
	assert.Equal(t, "myapp-configuration.yaml", cmList[0].Name)

	// secrets (index 1)
	secList, err := k.Client.Secret.List(namespace)
	assert.NoError(t, err)
	assert.Equal(t, "myapp-secret.yaml", secList[1].Name)

	// kratos release configuration (index 0)
	_, err = k.Client.Secret.Get(name+config.ConfigSuffix, namespace)
	assert.NoError(t, err)
}
