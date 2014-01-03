package httpcachetest

import (
	"net/http"
	"net/url"
)

func Request(path string) *http.Request {
	return &http.Request{URL: &url.URL{Path: path}}
}
