package hypermedia

import (
	"net/http"
	"net/url"
	"strings"
)

// TODO: need a full link header parser for http://tools.ietf.org/html/rfc5988
func HyperHeaderRelations(header http.Header, rels Relations) Relations {
	if rels == nil {
		rels = make(Relations)
	}

	link := header.Get("Link")
	if len(link) == 0 {
		return rels
	}

	for _, l := range strings.Split(link, ",") {
		l = strings.TrimSpace(l)
		segments := strings.Split(l, ";")

		if len(segments) < 2 {
			continue
		}

		if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
			continue
		}

		url, err := url.Parse(segments[0][1 : len(segments[0])-1])
		if err != nil {
			continue
		}

		link := Hyperlink(url.String())

		for _, segment := range segments[1:] {
			switch strings.TrimSpace(segment) {
			case `rel="next"`:
				rels["next"] = link
			case `rel="prev"`:
				rels["prev"] = link
			case `rel="first"`:
				rels["first"] = link
			case `rel="last"`:
				rels["last"] = link
			}
		}
	}

	return rels
}
