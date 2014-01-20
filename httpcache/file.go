package httpcache

import (
	"encoding/gob"
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
	relsFilename     = "rels"
	fileCreateFlag   = os.O_RDWR | os.O_CREATE | os.O_EXCL
)

// FileCache is a sawyer.Cacher that stores entries on disk.
type FileCache struct {
	path string
}

func NewFileCache(path string) *FileCache {
	return &FileCache{path}
}

func (c *FileCache) Get(req *http.Request) (sawyer.CachedResponse, error) {
	path := c.requestPath(req)

	responseFile, err := os.Open(filepath.Join(path, responseFilename))
	if err != nil {
		return nil, err
	}
	defer responseFile.Close()

	cachedResponse, err := Decode(responseFile)
	if err == nil {
		cachedResponse.Cacher = c
		cachedResponse.SetBodyFunc = func(res *sawyer.Response) {
			bodyFile, err := os.Open(filepath.Join(path, bodyFilename))
			if err == nil {
				res.Body = bodyFile
				res.BodyClosed = false
			} else {
				res.ResponseError = err
				res.BodyClosed = true
			}
		}
	}

	return cachedResponse, err
}

func (c *FileCache) Set(req *http.Request, res *sawyer.Response) error {
	path := c.requestPath(req)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	keyFile, err := newTempFile(path, keyFilename)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	keyFile.Write([]byte(RequestKey(req)))

	responseFile, err := newTempFile(path, responseFilename)
	if err != nil {
		return err
	}
	defer responseFile.Close()

	if err = Encode(res, responseFile); err != nil {
		return err
	}

	bodyFile, err := newTempFile(path, bodyFilename)
	if err != nil {
		return err
	}
	defer bodyFile.Close()

	err = EncodeBody(res, bodyFile)
	if err == nil {
		keyFile.Keep = true
		responseFile.Keep = true
		bodyFile.Keep = true
	}

	return err
}

func (c *FileCache) UpdateCache(req *http.Request, res *http.Response) error {
	path := c.requestPath(req)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	fullFile := filepath.Join(path, responseFilename)

	responseFile, err := os.Open(fullFile)
	if err != nil {
		return err
	}

	cached, err := Decode(responseFile)
	responseFile.Close()
	if err != nil {
		return err
	}

	cached.Expires = expiration(res)

	responseFile, err = os.Create(fullFile)
	if err != nil {
		return err
	}

	if err = EncodeResponse(cached.CachedResponse, responseFile); err != nil {
		return err
	}

	return nil
}

func (c *FileCache) SetRels(req *http.Request, rels hypermedia.Relations) error {
	path := c.requestPath(req)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	if len(rels) == 0 {
		os.Remove(filepath.Join(path, relsFilename))
		return nil
	}

	relsFile, err := newTempFile(path, relsFilename)
	if err != nil {
		return err
	}
	defer relsFile.Close()

	enc := gob.NewEncoder(relsFile)
	if err = enc.Encode(&rels); err == nil {
		relsFile.Keep = true
	}

	return err
}

func (c *FileCache) Rels(req *http.Request) (hypermedia.Relations, bool) {
	path := c.requestPath(req)
	relsFile, err := os.Open(filepath.Join(path, relsFilename))
	if err != nil {
		return nil, false
	}
	defer relsFile.Close()

	rels := make(hypermedia.Relations)
	dec := gob.NewDecoder(relsFile)
	dec.Decode(&rels)
	return rels, true
}

func (c *FileCache) requestPath(r *http.Request) string {
	sha := RequestSha(r)
	return filepath.Join(c.path, sha[0:2], sha[2:4], sha)
}

type tempFile struct {
	Name     string
	Temp     string
	Keep     bool
	tempFile *os.File
	file     *os.File
}

func newTempFile(dir string, name string) (*tempFile, error) {
	full := filepath.Join(dir, name)
	temp := filepath.Join(dir, "tmp_"+name)
	file, err := os.OpenFile(temp, fileCreateFlag, 0666)
	return &tempFile{full, temp, false, file, nil}, err
}

func (f *tempFile) Close() error {
	err := f.tempFile.Close()

	if f.Keep && err == nil {
		err = os.Rename(f.Temp, f.Name)
	}

	os.Remove(f.Temp)
	return err
}

func (f *tempFile) Write(data []byte) (int, error) {
	return f.tempFile.Write(data)
}
