package kratos

import (
	"bytes"
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
			k.PrintList(namespace)
		},
	)

	expected := "\x1b[39m\x1b[39m\x1b[96m\x1b[96mName \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mNamespace  \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mType      \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mReplicas\x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mCreation                     \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[96m\x1b[96mRevision\x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\n\x1b[39m\x1b[39mmyapp\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39mmynamespace\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39mdeployment\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m1       \x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m0001-01-01 00:00:00 +0000 UTC\x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39m0       \x1b[0m\x1b[0m\n"

	assert.Equal(t, expected, capture)
}
