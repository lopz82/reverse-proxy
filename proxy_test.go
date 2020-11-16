package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func Test_prepareRedirection(t *testing.T) {
	u1, _ := url.Parse("http://testing.com")
	testCases := []struct {
		dest url.URL
		r    *http.Request
	}{
		{*u1, &http.Request{}},
	}
	for _, test := range testCases {
		prepareRedirection(test.dest, test.r)
		assert.True(t, test.r.URL.Scheme == test.dest.Scheme)
		assert.True(t, test.r.URL.Host == test.dest.Host)
	}
}
