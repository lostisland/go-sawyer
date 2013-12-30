package httpcache

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetMissingCache(t *testing.T) {
	cache := NewMemoryCache()
	res := cache.Get("abc", nil)
	if res != nil {
		t.Fatal("response was found")
	}
}

func TestGetCacheWithoutValue(t *testing.T) {
	orig := &sawyer.Response{Response: &http.Response{StatusCode: 1}}

	cache := NewMemoryCache()
	cache.Set("abc", orig, nil)

	res := cache.Get("abc", nil)
	if res == nil {
		t.Fatal("Response is nil")
	}

	assert.Equal(t, false, res.IsError())
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

	cache := NewMemoryCache()
	err = cache.Set("abc", orig, testOrig)
	assert.Equal(t, nil, err)

	test := &TestResource{}
	res := cache.Get("abc", test)
	if res == nil {
		t.Fatal("Response is nil")
	}

	assert.Equal(t, false, res.IsError())
	assert.Equal(t, 1, res.StatusCode)
	assert.Equal(t, 2, test.A)
}

type TestResource struct {
	A int
}
