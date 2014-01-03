package httpcache

import (
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
	"os"
	"path/filepath"
)

const (
	keyFilename      = "key"
	responseFilename = "response"
	bodyFilename     = "body"
	fileCreateFlag   = os.O_RDWR | os.O_CREATE | os.O_EXCL
)

type FileCache struct {
	path string
}

func NewFileCache(path string) *FileCache {
	return &FileCache{path}
}

func (c *FileCache) Get(req *http.Request, v interface{}) *sawyer.Response {
	path := c.requestPath(req)

	responseFile, err := os.Open(filepath.Join(path, responseFilename))
	if err != nil {
		return ResponseError(err)
	}
	defer responseFile.Close()

	bodyFile, err := os.Open(filepath.Join(path, bodyFilename))
	if err != nil {
		return ResponseError(err)
	}
	defer bodyFile.Close()

	return DecodeFrom(v, responseFile, bodyFile)
}

func (c *FileCache) Set(req *http.Request, res *sawyer.Response, v interface{}) error {
	path := c.requestPath(req)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	keyFile, err := os.OpenFile(filepath.Join(path, keyFilename), fileCreateFlag, 0666)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	keyFile.Write([]byte(RequestKey(req)))

	responseFile, err := os.OpenFile(filepath.Join(path, responseFilename), fileCreateFlag, 0666)
	if err != nil {
		return err
	}
	defer responseFile.Close()

	bodyFile, err := os.OpenFile(filepath.Join(path, bodyFilename), fileCreateFlag, 0666)
	if err != nil {
		return err
	}
	defer bodyFile.Close()

	return EncodeTo(v, res, responseFile, bodyFile)
}

func (c *FileCache) Rels(req *http.Request) hypermedia.Relations {
	path := c.requestPath(req)

	responseFile, err := os.Create(filepath.Join(path, responseFilename))
	if err != nil {
		return hypermedia.Relations{}
	}
	defer responseFile.Close()

	return Decode(responseFile).Rels
}

func (c *FileCache) requestPath(r *http.Request) string {
	sha := RequestSha(r)
	return filepath.Join(c.path, sha[0:2], sha[2:4], sha)
}
