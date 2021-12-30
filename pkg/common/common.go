package common

import (
	"crypto/md5"
	"fmt"
)

// ListContain return true if searchItem is found in the list of string
func ListContain(list []string, searchItem string) bool {
	for _, item := range list {
		if item == searchItem {
			return true
		}
	}

	return false
}

// BoolPTR return a bool pointer
func BoolPTR(b bool) *bool {
	return &b
}

// MD5Sum return a md5sum from the input string
func MD5Sum(input string) string {
	hash := md5.New()
	return fmt.Sprintf("%x", hash.Sum([]byte(input)))
}
