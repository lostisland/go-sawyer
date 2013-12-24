package hypermedia

import (
	"net/url"
)

// HALResource is a resource with hypermedia specified as JSON HAL.
//
// http://stateless.co/hal_specification.html
type HALResource struct {
	Links Links `json:"_links"`
	rels  Relations
}

// Rels gets the link relations from the HALResource's Links field.
func (r *HALResource) Rels() Relations {
	if r.rels == nil {
		r.rels = make(map[string]Hyperlink)
		r.FillRels(r.rels)
	}
	return r.rels
}

func (r *HALResource) FillRels(rels Relations) {
	if r.Links == nil {
		return
	}

	for name, link := range r.Links {
		rels[name] = link.Href
	}
}

// Links is a collection of Link objects in a HALResource.  Note that the HAL
// spec allows single link objects or an array of link objects.  Sawyer
// currently only supports single link objects.
type Links map[string]Link

// Link represents a single link in a HALResource.
type Link struct {
	Href Hyperlink `json:"href"`
}

// Expand converts a uri template into a url.URL using the given M map.
func (l *Link) Expand(m M) (*url.URL, error) {
	return l.Href.Expand(m)
}
