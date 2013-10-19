package sawyer

import (
	"bytes"
	"encoding/json"
	"github.com/bmizerany/assert"
	"testing"
)

func TestHALRelations(t *testing.T) {
	input := `
{ "Login": "bob"
, "Url": "/foo/bar{/arg}"
, "_links":
	{ "self": { "href": "/self" }
	, "foo": { "href": "/foo" }
	, "bar": { "href": "/bar" }
	}
}`

	user := &HypermediaUser{}
	decode(t, input, user)

	rels := user.Rels()
	assert.Equal(t, 3, len(rels))
	assert.Equal(t, "/self", string(rels["self"]))
	assert.Equal(t, "/foo", string(rels["foo"]))
	assert.Equal(t, "/bar", string(rels["bar"]))

	rel, err := user.Rel("foo", nil)
	if err != nil {
		t.Fatalf("Error getting 'foo' relation: %s", err)
	}
	assert.Equal(t, "/foo", rel.Path)
}

func TestExpand(t *testing.T) {
	link := Hyperlink("/foo/bar{/arg}")
	u, _ := link.Expand(M{"arg": "baz", "foo": "bar"})
	assert.Equal(t, "/foo/bar/baz", u.String())
}

func TestExpandNil(t *testing.T) {
	link := Hyperlink("/foo/bar{/arg}")
	u, _ := link.Expand(nil)
	assert.Equal(t, "/foo/bar", u.String())
}

func TestDecode(t *testing.T) {
	input := `
{ "Login": "bob"
, "Url": "/foo/bar{/arg}"
, "_links":
  { "self": { "href": "/foo/bar{/arg}" }
  }
}`

	user := &HypermediaUser{}
	decode(t, input, user)

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

func decode(t *testing.T, input string, resource interface{}) {
	dec := json.NewDecoder(bytes.NewBufferString(input))
	err := dec.Decode(resource)
	if err != nil {
		t.Fatalf("Errors decoding json: %s", err)
	}
}

type HypermediaUser struct {
	Login string
	Url   Hyperlink
	*HALResource
}
