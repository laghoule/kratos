package kratos

import (
	"testing"

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
}
