package httpcache

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/httpcache/httpcachetest"
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMemoryGetMissingCache(t *testing.T) {
	req := httpcachetest.Request("abc")
	cache := NewMemoryCache()
	res := cache.Get(req)
	assert.Equal(t, true, res.IsError(), "response was found")
}

func TestMemoryGetCacheWithoutValue(t *testing.T) {
	orig := &sawyer.Response{Response: &http.Response{StatusCode: 1}}

	req := httpcachetest.Request("abc")
	cache := NewMemoryCache()
	cache.Set(req, orig)

	res := cache.Get(req)
	assert.Equal(t, false, res.IsError(), "response was not found")
	assert.Equal(t, 1, res.StatusCode)
}

func TestMemorySetAndGetCache(t *testing.T) {
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

	req := httpcachetest.Request("abc")
	cache := NewMemoryCache()
	err = cache.Set(req, orig)
	assert.Equal(t, nil, err)

	res := cache.Get(req)
	if res == nil {
		t.Fatal("Response is nil")
	}

	assert.Equal(t, false, res.IsError())
	assert.Equal(t, 1, res.StatusCode)

	test := &TestResource{}
	err = res.Decode(test)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, test.A)
}

type TestResource struct {
	A int
}
