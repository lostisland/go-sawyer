package httpcachetest

import (
	"github.com/lostisland/go-sawyer"
	"net/http"
	"net/url"
)

func Request(path string) *sawyer.Request {
	return &sawyer.Request{
		Request: &http.Request{URL: &url.URL{Path: path}},
	}
}
