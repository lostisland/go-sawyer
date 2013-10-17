package sawyer

import (
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
)

type Response struct {
	ResponseError error
	MediaType     *mediatype.MediaType
	isApiError    bool
	*http.Response
}

func (r *Request) Do(method string, output interface{}) *Response {
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

	res := &Response{nil, mtype, UseApiError(httpres.StatusCode), httpres}
	if mtype != nil {
		defer res.Body.Close()
		if res.isApiError {
			err = mtype.Decode(r.ApiError, res.Body)
		} else {
			err = mtype.Decode(output, res.Body)
		}
	}

	res.ResponseError = err
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

func ResponseError(err error) *Response {
	return &Response{ResponseError: err}
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
