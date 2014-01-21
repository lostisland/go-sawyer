package sawyer

import (
	"errors"
	"github.com/lostisland/go-sawyer/hypermedia"
	"github.com/lostisland/go-sawyer/mediatype"
	"io"
	"net/http"
)

// Response is a wrapped net/http Response with a pointer to the MediaType and
// the cacher.  It also doubles as a possible error object.
type Response struct {
	// ResponseError stores any errors made making the HTTP request.  If set, then
	// AnyError() and IsError() will return true, and Error() will delegate to it.
	ResponseError error
	MediaType     *mediatype.MediaType
	BodyClosed    bool
	Cacher        Cacher
	isApiError    bool
	rels          hypermedia.Relations
	*http.Response
}

// AnyError returns true if the HTTP request returned an error, or if the
// response status is not a 2xx code.
func (r *Response) AnyError() bool {
	return r.IsError() || r.IsApiError()
}

// IsError returns true if the HTTP request returned an error.
func (r *Response) IsError() bool {
	return r.ResponseError != nil
}

// IsApiError returns true if the response status is not a 2xx code.
func (r *Response) IsApiError() bool {
	return r.isApiError
}

// Error returns the ResponseError's error string if set, or an empty string.
func (r *Response) Error() string {
	if r.ResponseError != nil {
		return r.ResponseError.Error()
	}
	return ""
}

// Decode will decode the body into the given resource, and parse the hypermedia
// relations.  This is meant to be called after an HTTP request, and will close
// the response body.  The decoder is set from the response's MediaType.
func (r *Response) Decode(resource interface{}) error {
	if r.ResponseError != nil || r.Response == nil {
		return r.ResponseError
	}

	if r.MediaType == nil {
		return errorNoMediaType
	}

	if resource == nil {
		return errorNoResource
	}

	if r.BodyClosed {
		return errorBodyClosed
	}

	defer r.Body.Close()
	r.BodyClosed = true

	if err := r.DecodeFrom(resource, r.Body); err != nil {
		return err
	}

	rels := hypermedia.Rels(resource)
	if err := r.Cacher.SetRels(r.Request, rels); err != nil {
		return err
	}

	if cachedResource, ok := resource.(hypermedia.CachedResource); ok {
		cachedResource.CacheRels(rels)
	}

	return nil
}

// DecodeFrom decodes the resource from the given io.Reader, using the decoder
// from the response's MediaType.
func (r *Response) DecodeFrom(resource interface{}, body io.Reader) error {
	if resource == nil {
		return errorNoResource
	}

	if r.ContentLength < 1 {
		return errorNoBody
	}

	dec, err := r.MediaType.Decoder(body)
	if err != nil {
		return err
	}

	if err := dec.Decode(resource); err != nil {
		return err
	}

	return nil
}

// HypermediaRels implements the hypermedia.HypermediaResource interface.  The
// relations are parsed from the Link header.
func (r *Response) HypermediaRels(rels hypermedia.Relations) {
	if r.rels == nil {
		r.rels = hypermedia.HyperHeaderRelations(r.Header, nil)
	}

	for key, value := range r.rels {
		rels[key] = value
	}
}

// Rels returns the cached relations if they have been set.
func (r *Response) Rels() (hypermedia.Relations, bool) {
	if r.rels != nil {
		return r.rels, true
	}
	return nil, false
}

// CacheRels will set the given relations for this resource.
func (r *Response) CacheRels(rels hypermedia.Relations) {
	r.rels = rels
}

// NewResponse initializes a new Response with common internal values set.
func NewResponse(res *http.Response) *Response {
	if res == nil {
		res = &http.Response{}
	}

	return &Response{
		Response:   res,
		Cacher:     noOpCacher,
		BodyClosed: false,
		isApiError: UseApiError(res.StatusCode),
	}
}

// ResponseError returns an empty Response with the ResponseError set from the
// given error.
func ResponseError(err error) *Response {
	res := NewResponse(nil)
	res.ResponseError = err
	res.BodyClosed = true
	return res
}

// UseApiError determines if the given status is considered an API error.
func UseApiError(status int) bool {
	switch {
	case status > 199 && status < 300:
		return false
	case status == 304:
		return false
	case status == 0:
		return false
	}
	return true
}

func mediaType(res *http.Response) (*mediatype.MediaType, error) {
	if ctype := res.Header.Get(ctypeHeader); len(ctype) > 0 {
		return mediatype.Parse(ctype)
	}
	return nil, nil
}

var (
	errorNoMediaType = errors.New("No media type for this response")
	errorNoResource  = errors.New("No resource value provided")
	errorNoBody      = errors.New("No response body to decode")
	errorBodyClosed  = errors.New("Response body has already been read")
)
