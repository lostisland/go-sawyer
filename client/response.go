package sawyer

import (
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
)

type Response struct {
	MediaType  *mediatype.MediaType
	isApiError bool
	*http.Response
}

func (r *Request) Get(output interface{}) (*Response, error) {
	return r.Do(GetMethod, output)
}

func (r *Request) Do(method string, output interface{}) (*Response, error) {
	httpres, err := r.Client.Do(r.Request)
	if err != nil {
		return nil, err
	}

	mtype, err := mediaType(httpres)
	if err != nil {
		httpres.Body.Close()
		return nil, err
	}

	res := &Response{mtype, UseApiError(httpres.StatusCode), httpres}
	if mtype != nil {
		defer res.Body.Close()
		if res.isApiError {
			err = mtype.Decode(r.ApiError, res.Body)
		} else {
			err = mtype.Decode(output, res.Body)
		}
	}

	return res, nil
}

func (r *Response) IsApiError() bool {
	return r.isApiError
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
	if ctype := res.Header.Get("Content-Type"); len(ctype) > 0 {
		return mediatype.Parse(ctype)
	}
	return nil, nil
}
