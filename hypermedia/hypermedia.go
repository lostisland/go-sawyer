package hypermedia

import (
	"fmt"
	"github.com/jtacoma/uritemplates"
	"net/url"
	"reflect"
)

type Links map[string]Link

type Link struct {
	Href Hyperlink `json:"href"`
}

type Hyperlink string

func (l *Hyperlink) Expand(m map[string]interface{}) (*url.URL, error) {
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

func (l *Link) Expand(m map[string]interface{}) (*url.URL, error) {
	return l.Href.Expand(m)
}

type Relations map[string]Hyperlink

func (h Relations) Rel(name string, m map[string]interface{}) (*url.URL, error) {
	if rel, ok := h[name]; ok {
		return rel.Expand(m)
	}
	return nil, fmt.Errorf("No %s relation found", name)
}

type HypermediaResource interface {
	Rels() Relations
}

// HALResource is a resource with hypermedia specified as JSON HAL.
type HALResource struct {
	Links Links `json:"_links"`
	rels  Relations
}

func (r *HALResource) Rels() Relations {
	if r.rels == nil {
		r.rels = make(map[string]Hyperlink)
		for name, link := range r.Links {
			r.rels[name] = link.Href
		}
	}
	return r.rels
}

func HypermediaDecoder(res HypermediaResource) Relations {
	return res.Rels()
}

func HyperFieldDecoder(res interface{}) Relations {
	rels := make(Relations)
	t := reflect.TypeOf(res).Elem()
	v := reflect.ValueOf(res).Elem()
	fieldlen := t.NumField()
	for i := 0; i < fieldlen; i++ {
		fillRelation(rels, t, v, i)
	}
	return rels
}

func fillRelation(rels map[string]Hyperlink, t reflect.Type, v reflect.Value, index int) {
	f := t.Field(index)

	if hyperlinkType != f.Type {
		return
	}

	hl := v.Field(index).Interface().(Hyperlink)
	name := f.Name
	if rel := f.Tag.Get("rel"); len(rel) > 0 {
		name = rel
	}
	rels[name] = hl
}

var hyperlinkType = reflect.TypeOf(Hyperlink("foo"))
