// Package hypermedia provides helpers for parsing hypermedia links in resources
// and expanding the links to make further requests.
package hypermedia

import (
	"fmt"
	"github.com/jtacoma/uritemplates"
	"net/url"
)

// Relations is a map of keys that point to Hyperlink objects.
type Relations map[string]Hyperlink

// Rels returns a new Relations object.
func NewRels() Relations {
	return Relations{}
}

// Rel fetches and expands the Hyperlink by its given key in the Relations map.
func (h Relations) Rel(name string, m M) (*url.URL, error) {
	if rel, ok := h[name]; ok {
		return rel.Expand(m)
	}
	return nil, fmt.Errorf("No %s relation found", name)
}

// Rels gets the hypermedia relations from the given resource.
func Rels(resource interface{}) Relations {
	cachedResource, ok := resource.(CachedResource)
	if ok {
		if rels, cached := cachedResource.Rels(); cached {
			return rels
		}
	}

	rels := NewRels()
	FillRels(resource, rels)

	if ok {
		cachedResource.CacheRels(rels)
	}

	return rels
}

// FillRels populates the given relations object from the relations in the
// resource.
func FillRels(resource interface{}, rels Relations) {
	if hypermediaRel, ok := resource.(HyperfieldResource); ok {
		HyperFieldRelations(hypermediaRel, rels)
	}

	if hypermediaRel, ok := resource.(HypermediaResource); ok {
		hypermediaRel.HypermediaRels(rels)
	}
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

// A HypermediaResource has link relations for next actions of a resource.
type HypermediaResource interface {
	HypermediaRels(Relations)
}

// A CachedResource is capable of caching the relations locally, so that
// multiple accesses don't require parsing it again.
type CachedResource interface {
	Rels() (Relations, bool)
	CacheRels(Relations)
}
