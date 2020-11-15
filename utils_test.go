package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateUrls(t *testing.T) {

	returnsError := []string{
		"http://testing.com",
		"http://testing2.com",
		":/testing2.com",
	}
	_, err := validateURLs(returnsError)
	assert.Error(t, err)

	valid := []string{
		"http://testing.com",
		"http://testing2.com",
	}
	servers, err := validateURLs(valid)
	assert.NoError(t, err)
	assert.Len(t, servers, 2)

	for i := range servers {
		assert.True(t, valid[i] == servers[i].String())
	}
}
