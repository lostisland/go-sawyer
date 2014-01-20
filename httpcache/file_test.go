package httpcache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFile(t *testing.T) {
	setup := FileSetup(t)
	defer setup.Teardown()
	CacheResponsesTestFor(setup.Cache, t)
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
