package sawyer

import (
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
)

type Response struct {
	ResponseError error
	MediaType     *mediatype.MediaType
	isApiError    bool
	BodyClosed    bool
	*http.Response
}

func (r *Request) Do(method string, output interface{}) *Response {
	r.URL.RawQuery = r.Query.Encode()
	r.Method = method
	httpres, err := r.Client.Do(r.Request)
	if err != nil {
		return ResponseError(err)
	}

	mtype, err := mediaType(httpres)
	if err != nil {
		httpres.Body.Close()
		return ResponseError(err)
	}

	res := &Response{nil, mtype, UseApiError(httpres.StatusCode), false, httpres}
	if mtype != nil {
		res.decode(r.ApiError, output)
	}

	return res
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

func (r *Response) decode(apierr interface{}, output interface{}) {
	if r.isApiError {
		r.decodeResource(apierr)
	} else {
		r.decodeResource(output)
	}
}

func (r *Response) decodeResource(resource interface{}) {
	if resource == nil {
		return
	}

	dec := r.MediaType.Decoder(r.Body)
	if dec == nil {
		return
	}

	defer r.Body.Close()
	r.BodyClosed = true
	r.ResponseError = dec.Decode(resource)
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
