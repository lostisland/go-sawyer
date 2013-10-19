package sawyer

import (
	"fmt"
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
)

type Response struct {
	ResponseError error
	MediaType     *mediatype.MediaType
	isApiError    bool
	ApiError      interface{}
	BodyClosed    bool
	errorFunc     func() interface{}
	*http.Response
}

func (r *Response) AnyError() bool {
	return r.IsError() || r.IsApiError()
}

func (r *Response) IsError() bool {
	return r.ResponseError != nil
}

func (r *Response) IsApiError() bool {
	return r.ApiError != nil
}

func (r *Response) Error() string {
	if r.ResponseError != nil {
		return r.ResponseError.Error()
	}
	return ""
}

func (r *Response) decode(output interface{}) {
	if r.isApiError {
		r.ApiError = r.errorFunc()
		r.decodeResource(r.ApiError)
	} else {
		r.decodeResource(output)
	}
}

func (r *Response) decodeResource(resource interface{}) {
	if resource == nil {
		return
	}

	defer r.Body.Close()
	r.BodyClosed = true

	dec := r.MediaType.Decoder(r.Body)
	if dec == nil {
		r.ResponseError = fmt.Errorf("No decoder found for format %s (%s)",
			r.MediaType.Format, r.MediaType.String())
	} else {
		r.ResponseError = dec.Decode(resource)
	}
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
