package kratos

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/pterm/pterm"
	"github.com/stretchr/testify/assert"
)

func TestIsReleaseExist(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.Create(name, namespace); err != nil {
		t.Error(err)
		return
	}

	notfound, err := k.IsReleaseExist(name, namespace)
	if err != nil {
		t.Error(err)
		return
	}

	assert.False(t, notfound)
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

	assert.Len(t, list, 0)
}

func listOutput(f func()) string {
	var buf bytes.Buffer
	pterm.SetDefaultOutput(&buf)

	f()

	pterm.SetDefaultOutput(os.Stderr)
	return buf.String()
}

func TestPrintList(t *testing.T) {
	k := new()
	k.Config = createConf()

	if err := k.Create(name, namespace); err != nil {
		t.Error(err)
		return
	}

	capture := listOutput(
		func() {
			if err := k.PrintList(namespace); err != nil {
				fmt.Println(err)
				return
			}
		},
	)

	expected := "\x1b[39m\x1b[39m\x1b[96m\x1b[96mName\x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mNamespace\x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mType\x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mReplicas\x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mCreation\x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mRevision\x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\n"

	assert.Equal(t, expected, capture)
}
