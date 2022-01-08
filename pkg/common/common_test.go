package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolPTR(t *testing.T) {
	expected := true
	assert.Equal(t, &expected, BoolPTR(true))
}

func TestMD5Sum16(t *testing.T) {
	expected := "74657374d41d8cd9"
	assert.Equal(t, expected, MD5Sum16("test"))
}

func TestListContain(t *testing.T) {
	list := []string{"one", "two", "three"}
	found := ListContain(list, "two"); 
	assert.True(t, found)

	found = ListContain(list, "four")
	assert.False(t, found)
}
