package sawyer

import (
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
