package memorycache

import (
	"bytes"
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/httpcache"
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
	"testing"
)

func TestGetMissingCache(t *testing.T) {
	cache := New()
	res := cache.Get("abc", nil)
	if res != nil {
		t.Fatal("response was found")
	}
}

func TestGetCacheWithoutValue(t *testing.T) {
	orig := &sawyer.Response{Response: &http.Response{StatusCode: 1}}

	var buf bytes.Buffer
	err := httpcache.Encode(orig, &buf)
	assert.Equal(t, nil, err)

	cache := New()
	cache.Cache["abc"] = &cacheEntry{&buf, nil}

	res := cache.Get("abc", nil)
	if res == nil {
		t.Fatal("Response is nil")
	}

	assert.Equal(t, false, res.IsError())
	assert.Equal(t, 1, res.StatusCode)
}

func TestGetCacheWithEmptyValue(t *testing.T) {
	orig := &sawyer.Response{Response: &http.Response{StatusCode: 1}}

	var buf bytes.Buffer
	err := httpcache.Encode(orig, &buf)
	assert.Equal(t, nil, err)

	cache := New()
	cache.Cache["abc"] = &cacheEntry{&buf, nil}

	test := &TestResource{}
	res := cache.Get("abc", test)
	if res == nil {
		t.Fatal("Response is nil")
	}

	assert.Equal(t, false, res.IsError())
	assert.Equal(t, 1, res.StatusCode)
	assert.Equal(t, 0, test.A)
}

func TestGetCacheWithValue(t *testing.T) {
	orig := &sawyer.Response{Response: &http.Response{StatusCode: 1}}
	mt, err := mediatype.Parse("application/json")
	assert.Equal(t, nil, err)
	orig.MediaType = mt

	testOrig := &TestResource{2}
	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	err = enc.Encode(testOrig)
	assert.Equal(t, nil, err)

	orig.Response.ContentLength = int64(body.Len())

	var buf bytes.Buffer
	err = httpcache.Encode(orig, &buf)
	assert.Equal(t, nil, err)

	cache := New()
	cache.Cache["abc"] = &cacheEntry{&buf, &body}

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
