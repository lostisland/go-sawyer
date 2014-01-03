package httpcache

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/httpcache/httpcachetest"
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestFileGetMissingCache(t *testing.T) {
	setup := FileSetup(t)
	defer setup.Teardown()

	req := httpcachetest.Request("abc")
	cache := setup.Cache
	res := cache.Get(req, nil)
	assert.Equal(t, true, res.IsError())
}

func TestFileSetAndGetCache(t *testing.T) {
	setup := FileSetup(t)
	defer setup.Teardown()

	mt, err := mediatype.Parse("application/json")
	assert.Equal(t, nil, err)

	testOrig := &FileTestResource{2}
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
	cache := setup.Cache
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

type fileSetup struct {
	Path  string
	Cache *FileCache
	*testing.T
}

func FileSetup(t *testing.T) *fileSetup {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(wd, "filecachetest")
	err = os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatal(err)
	}

	return &fileSetup{path, NewFileCache(path), t}
}

func (s *fileSetup) Teardown() {
	if err := os.RemoveAll(s.Path); err != nil {
		s.Fatal(err)
	}
}

type FileTestResource struct {
	A int
}
