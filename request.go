package sawyer

import (
	"github.com/lostisland/go-sawyer/hypermedia"
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	Client    *http.Client
	MediaType *mediatype.MediaType
	Query     url.Values
	cacher    Cacher
	*http.Request
}

// NewRequest creates a new sawyer.Request for the given relative url path, with
// any default headers or query parameters specified on Client.  The Request URL
// is resolved to an absolute URL.
func (c *Client) NewRequest(rawurl string) (*Request, error) {
	httpreq, err := buildRequest(c, rawurl)
	if httpreq == nil {
		return nil, err
	}

	return &Request{c.HttpClient, nil, httpreq.URL.Query(), c.Cacher, httpreq}, err
}

func (r *Request) Do(method string) *Response {
	r.URL.RawQuery = r.Query.Encode()
	r.Method = method
	cached := r.cacher.Get(r.Request)
	if !cached.IsError() {
		return cached
	}

	httpres, err := r.Client.Do(r.Request)
	if err != nil {
		return ResponseError(err)
	}

	mtype, err := mediaType(httpres)
	if err != nil {
		httpres.Body.Close()
		return ResponseError(err)
	}

	res := &Response{
		MediaType:  mtype,
		BodyClosed: false,
		Response:   httpres,
		rels:       hypermedia.HyperHeaderRelations(httpres.Header, nil),
		isApiError: UseApiError(httpres.StatusCode),
	}

	if !res.AnyError() {
		r.cacher.Set(r.Request, res)
	}

	return res
}

func (r *Request) Head() *Response {
	return r.Do(HeadMethod)
}

func (r *Request) Get() *Response {
	return r.Do(GetMethod)
}

func (r *Request) Post() *Response {
	return r.Do(PostMethod)
}

func (r *Request) Put() *Response {
	return r.Do(PutMethod)
}

func (r *Request) Patch() *Response {
	return r.Do(PatchMethod)
}

func (r *Request) Delete() *Response {
	return r.Do(DeleteMethod)
}

func (r *Request) Options() *Response {
	return r.Do(OptionsMethod)
}

// Encodes and sets the proper headers for the request body.
func (r *Request) SetBody(mtype *mediatype.MediaType, input interface{}) error {
	r.MediaType = mtype
	buf, err := mtype.Encode(input)
	if err != nil {
		return err
	}

	r.Header.Set(ctypeHeader, mtype.String())
	r.ContentLength = int64(buf.Len())
	r.Body = ioutil.NopCloser(buf)
	return nil
}

const (
	ctypeHeader   = "Content-Type"
	HeadMethod    = "HEAD"
	GetMethod     = "GET"
	PostMethod    = "POST"
	PutMethod     = "PUT"
	PatchMethod   = "PATCH"
	DeleteMethod  = "DELETE"
	OptionsMethod = "OPTIONS"
)
