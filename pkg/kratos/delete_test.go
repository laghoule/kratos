package kratos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.Create(name, namespace); err != nil {
		t.Error(err)
		return
	}

	svcList, err := k.Client.Service.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEmpty(t, svcList)

	cmList, err := k.Client.ConfigMaps.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEmpty(t, cmList)

	secList, err := k.Client.Secrets.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEmpty(t, secList)

	if err := k.Delete(name, namespace); err != nil {
		t.Error(err)
		return
	}

	svcList, err = k.Client.Service.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Empty(t, svcList)

	cmList, err = k.Client.ConfigMaps.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Empty(t, cmList)

	secList, err = k.Client.Secrets.List(namespace)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Empty(t, secList)
}
