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

	// service
	svcList, err := k.Client.Service.List(namespace)
	assert.NoError(t, err)
	assert.Equal(t, name, svcList[0].Name)

	// secrets (index 1)
	secList, err := k.Client.Secrets.List(namespace)
	assert.NoError(t, err)
	assert.Equal(t, "myapp-secret.yaml", secList[1].Name)

	// kratos release configuration (index 0)
	_, err = k.Client.Secrets.Get(name+config.ConfigSuffix, namespace)
	assert.NoError(t, err)
}
