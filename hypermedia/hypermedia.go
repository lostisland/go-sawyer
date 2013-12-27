// Package hypermedia provides helpers for parsing hypermedia links in resources
// and expanding the links to make further requests.
package hypermedia

import (
	"fmt"
	"github.com/jtacoma/uritemplates"
	"net/url"
)

func Rels() Relations {
	return Relations{}
}

// Hyperlink is a string url.  If it is a uri template, it can be converted to
// a full URL with Expand().
type Hyperlink string

// Expand converts a uri template into a url.URL using the given M map.
func (l Hyperlink) Expand(m M) (*url.URL, error) {
	template, err := uritemplates.Parse(string(l))
	if err != nil {
		return nil, err
	}

	// clone M to map[string]interface{}
	// if we don't do this type assertion will
	// fail on jtacoma/uritemplates
	// see https://github.com/jtacoma/uritemplates/blob/master/uritemplates.go#L189
	mm := make(map[string]interface{}, len(m))
	for k, v := range m {
		mm[k] = v
	}

	expanded, err := template.Expand(mm)
	if err != nil {
		return nil, err
	}

	return url.Parse(expanded)
}

// M represents a map of values to expand a Hyperlink.
type M map[string]interface{}

// Relations is a map of keys that point to Hyperlink objects.
type Relations map[string]Hyperlink

// Rel fetches and expands the Hyperlink by its given key in the Relations map.
func (h Relations) Rel(name string, m M) (*url.URL, error) {
	if rel, ok := h[name]; ok {
		return rel.Expand(m)
	}
	return nil, fmt.Errorf("No %s relation found", name)
}

// A HypermediaResource has link relations for next actions of a resource.
type HypermediaResource interface {
	FillRels(Relations)
}
