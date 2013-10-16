package hypermedia

import (
	"bytes"
	"encoding/json"
	"github.com/bmizerany/assert"
	"testing"
)

func TestDecode(t *testing.T) {
	input := `
{ "Login": "bob"
, "Url": "/foo/bar{/arg}"
, "_links":
  { "self": { "href": "/foo/bar{/arg}" }
  }
}`

	user := &TestUser{}
	dec := json.NewDecoder(bytes.NewBufferString(input))
	err := dec.Decode(user)
	if err != nil {
		t.Fatalf("Errors decoding json: %s", err)
	}

	assert.Equal(t, "bob", user.Login)
	assert.Equal(t, 1, len(user.Links))

	hl := user.Url
	url, err := hl.Expand(M{"arg": "baz"})
	if err != nil {
		t.Errorf("Errors parsing %s: %s", hl, err)
	}

	assert.Equal(t, "/foo/bar/baz", url.String())

	hl = user.Links["self"].Href
	url, err = hl.Expand(M{"arg": "baz"})
	if err != nil {
		t.Errorf("Errors parsing %s: %s", hl, err)
	}
	assert.Equal(t, "/foo/bar/baz", url.String())
}

type TestUser struct {
	Login string
	Url   Hyperlink
	Links Links `json:"_links"`
}
