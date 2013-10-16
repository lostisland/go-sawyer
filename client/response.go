package sawyer

import (
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
)

type Response struct {
	ResponseError error
	MediaType     *mediatype.MediaType
	useApiError   bool
	*http.Response
}

func (r *Response) IsError() bool {
	return r.ResponseError != nil
}

func (r *Response) IsApiError() bool {
	return r.useApiError
}

func (r *Response) Error() string {
	if r.IsError() {
		return r.ResponseError.Error()
	}
	return ""
}

func apiResponse(err error) *Response {
	return &Response{ResponseError: err}
}

func buildResponse(resource interface{}, apierr interface{}, c *Client, httpres *http.Response, err error) *Response {
	res := &Response{ResponseError: err, Response: httpres}

	if err != nil {
		return res
	}

	defer httpres.Body.Close()

	mt, err := mediaType(httpres)
	if err != nil {
		res.ResponseError = err
		return res
	}

	res.MediaType = mt
	res.useApiError = UseApiError(httpres.StatusCode)
	res.ResponseError = decode(resource, apierr, c, res)
	return res
}

func decode(resource interface{}, apierr interface{}, c *Client, res *Response) error {
	decoder := res.MediaType.Decoder(res.Body)
	if decoder == nil {
		return nil
	}

	if res.useApiError && apierr != nil {
		return decoder.Decode(apierr)
	} else if resource != nil {
		return decoder.Decode(resource)
	}
	return nil
}

func mediaType(res *http.Response) (*mediatype.MediaType, error) {
	if ctype := res.Header.Get("Content-Type"); len(ctype) > 0 {
		return mediatype.Parse(ctype)
	}
	return nil, nil
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
