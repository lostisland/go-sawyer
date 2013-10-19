package sawyer

import (
	"fmt"
	"github.com/jtacoma/uritemplates"
	"net/url"
)

type Links map[string]Link

type Link struct {
	Href Hyperlink `json:"href"`
}

type Hyperlink string

// uri template map
type M map[string]interface{}

func (l *Hyperlink) Expand(m M) (*url.URL, error) {
	template, err := uritemplates.Parse(string(*l))
	if err != nil {
		return nil, err
	}

	expanded, err := template.Expand(m)
	if err != nil {
		return nil, err
	}

	return url.ParseRequestURI(expanded)
}

func (l *Link) Expand(m M) (*url.URL, error) {
	return l.Href.Expand(m)
}

// HypermediaResource describes any REST resource with hypermedia relations.
type HypermediaResource interface {
	Rels() map[string]Hyperlink
	Rel(string, M) (*url.URL, error)
}

// HALResource is a resource with hypermedia specified as JSON HAL.
type HALResource struct {
	Links Links `json:"_links"`
	rels  map[string]Hyperlink
}

func (r *HALResource) Rels() map[string]Hyperlink {
	if r.rels == nil {
		r.rels = make(map[string]Hyperlink)
		for name, link := range r.Links {
			r.rels[name] = link.Href
		}
	}
	return r.rels
}

func (r *HALResource) Rel(name string, m M) (*url.URL, error) {
	if rel, ok := r.Rels()[name]; ok {
		return rel.Expand(m)
	}
	return nil, fmt.Errorf("No %s relation found", name)
}
