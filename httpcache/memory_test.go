package httpcache

import (
	"testing"
)

func TestMemory(t *testing.T) {
	CacheResponsesTestFor(NewMemoryCache(), t)
}
