package httpcache

import (
	"github.com/lostisland/go-sawyer"
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

func (c *FileCache) Get(req *http.Request) *sawyer.Response {
	path := c.requestPath(req)

	responseFile, err := os.Open(filepath.Join(path, responseFilename))
	if err != nil {
		return ResponseError(err)
	}
	defer responseFile.Close()

	res := Decode(responseFile)
	res.Request = req

	bodyFile, err := os.Open(filepath.Join(path, bodyFilename))
	if err != nil {
		res.ResponseError = err
		res.BodyClosed = true
		return res
	}

	res.Body = bodyFile
	return res
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
