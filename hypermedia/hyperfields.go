package hypermedia

import (
	"reflect"
)

type HyperfieldResource interface {
	HyperfieldRels()
}

// The HyperFieldRelations gets link relations from a resource by reflecting on
// its Hyperlink properties.  The relation name is taken either from the name
// of the field, or a "rel" struct tag.
//
//   type Foo struct {
//     Url         Hyperlink `rel:"self" json:"url"`
//     CommentsUrl Hyperlink `rel:"comments" json:"comments_url"`
//   }
//
func HyperFieldRelations(res interface{}, rels Relations) Relations {
	if rels == nil {
		rels = make(Relations)
	}
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
