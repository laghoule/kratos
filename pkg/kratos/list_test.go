package kratos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsReleaseExist(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.Create(name, namespace); err != nil {
		t.Error(err)
		return
	}

	found, err := k.IsReleaseExist(name, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.True(t, found)
}

func TestGetList(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.Create(name, namespace); err != nil {
		t.Error(err)
		return
	}

	list, err := k.GetList(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, list, 1)
}

func TestPrintList(t *testing.T) {
	// TODO TestPrintList
}
