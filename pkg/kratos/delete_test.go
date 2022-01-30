package kratos

import (
	"testing"
)

func TestDelete(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.Delete(name, namespace); err != nil {
		t.Error(err)
		return
	}
}
