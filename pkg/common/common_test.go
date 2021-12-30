package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolPTR(t *testing.T) {
	expected := true
	assert.Equal(t, &expected, BoolPTR(true))
}

func TestMD5sum(t *testing.T) {
	expected := "74657374d41d8cd98f00b204e9800998ecf8427e"
	assert.Equal(t, expected, MD5Sum("test"))
}

func TestListContain(t *testing.T) {
	// TODO: TestListContain
}
