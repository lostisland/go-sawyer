package sawyer

import (
	"errors"
	"github.com/lostisland/go-sawyer/hypermedia"
	"github.com/lostisland/go-sawyer/mediatype"
	"io"
	"net/http"
)

type Response struct {
	ResponseError error
	MediaType     *mediatype.MediaType
	BodyClosed    bool
	isApiError    bool
	rels          hypermedia.Relations
	*http.Response
}

func (r *Response) AnyError() bool {
	return r.IsError() || r.IsApiError()
}

func (r *Response) IsError() bool {
	return r.ResponseError != nil
}

func (r *Response) IsApiError() bool {
	return r.isApiError
}

func (r *Response) Error() string {
	if r.ResponseError != nil {
		return r.ResponseError.Error()
	}
	return ""
}

func (r *Response) HypermediaRels(rels hypermedia.Relations) {
	for key, value := range r.rels {
		rels[key] = value
	}
}

func (r *Response) Decode(resource interface{}) error {
	if r.MediaType == nil {
		return errors.New("No media type for this response")
	}

	if resource == nil || r.ResponseError != nil || r.BodyClosed {
		return r.ResponseError
	}

	defer r.Body.Close()
	r.BodyClosed = true

	r.ResponseError = r.DecodeFrom(resource, r.Body)

	return r.ResponseError
}

func (r *Response) DecodeFrom(resource interface{}, body io.Reader) error {
	if resource == nil || r.ContentLength < 1 {
		return nil
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

func ResponseError(err error) *Response {
	return &Response{ResponseError: err, BodyClosed: true}
}

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
