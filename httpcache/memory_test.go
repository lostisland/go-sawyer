package httpcache

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestGetMissingCache(t *testing.T) {
	req := request()
	cache := NewMemoryCache()
	res := cache.Get(req, nil)
	assert.Equal(t, true, res.IsError(), "response was found")
}

func TestGetCacheWithoutValue(t *testing.T) {
	orig := &sawyer.Response{Response: &http.Response{StatusCode: 1}}

	req := request()
	cache := NewMemoryCache()
	cache.Set(req, orig, nil)

	res := cache.Get(req, nil)
	assert.Equal(t, false, res.IsError(), "response was not found")
	assert.Equal(t, 1, res.StatusCode)
}

func TestSetAndGetCache(t *testing.T) {
	mt, err := mediatype.Parse("application/json")
	assert.Equal(t, nil, err)

	testOrig := &TestResource{2}
	body, err := mt.Encode(testOrig)

	orig := &sawyer.Response{
		MediaType: mt,
		Response: &http.Response{
			StatusCode:    1,
			ContentLength: int64(body.Len()),
			Body:          ioutil.NopCloser(body),
		},
	}

	req := request()
	cache := NewMemoryCache()
	err = cache.Set(req, orig, testOrig)
	assert.Equal(t, nil, err)

	test := &TestResource{}
	res := cache.Get(req, test)
	if res == nil {
		t.Fatal("Response is nil")
	}

	assert.Equal(t, false, res.IsError())
	assert.Equal(t, 1, res.StatusCode)
	assert.Equal(t, 2, test.A)
}

func request() *http.Request {
	return &http.Request{URL: &url.URL{Path: "abc"}}
}

type TestResource struct {
	A int
}
