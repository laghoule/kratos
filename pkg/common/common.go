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

// PTR return a pointer
func PTR[T any](t T) *T {
	return &t
}

// MD5Sum16 return a md5sum cut to 16 caracters
func MD5Sum16(input string) string {
	hash := md5.New()
	md5sum := fmt.Sprintf("%x", hash.Sum([]byte(input)))

	if len(md5sum) > 16 {
		return md5sum[0:16]
	}

	return md5sum
}
