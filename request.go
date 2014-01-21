package sawyer

import (
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Request is a wrapped net/http Request with a pointer to the net/http Client,
// MediaType, parsed URI query, and the configured Cacher.  Requests are capable
// of returning a sawyer Response with Do() or the HTTP verb helpers (Get(),
// Head(), Post(), etc).
type Request struct {
	Client    *http.Client
	MediaType *mediatype.MediaType
	Query     url.Values
	Cacher    Cacher
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

// Do completes the HTTP request, returning a response.  The Request's Cacher is
// used to return a cached response if available.  Otherwise, the request goes
// through and fills the cache for future requests.
func (r *Request) Do(method string) *Response {
	r.URL.RawQuery = r.Query.Encode()
	r.Method = method
	cached, cachedErr := r.Cacher.Get(r.Request)
	if cachedErr == nil {
		if cached.IsFresh() {
			return cached.Decode(r)
		} else {
			cached.SetupRequest(r.Request)
		}
	}

	httpres, err := r.Client.Do(r.Request)
	if err != nil {
		return ResponseError(err)
	}

	if cachedErr == nil && httpres.StatusCode == 304 {
		r.Cacher.UpdateCache(r.Request, httpres)
		return cached.Decode(r)
	}

	mtype, err := mediaType(httpres)
	if err != nil {
		httpres.Body.Close()
		return ResponseError(err)
	}

	res := NewResponse(httpres)
	res.MediaType = mtype
	res.Cacher = r.Cacher

	if !res.AnyError() {
		r.Cacher.Set(r.Request, res)
	}

	return res
}

// Head is a helper method for Do().
func (r *Request) Head() *Response {
	return r.Do(HeadMethod)
}

// Get is a helper method for Do().
func (r *Request) Get() *Response {
	return r.Do(GetMethod)
}

// Post is a helper method for Do().
func (r *Request) Post() *Response {
	return r.Do(PostMethod)
}

// Put is a helper method for Do().
func (r *Request) Put() *Response {
	return r.Do(PutMethod)
}

// Patch is a helper method for Do().
func (r *Request) Patch() *Response {
	return r.Do(PatchMethod)
}

// Delete is a helper method for Do().
func (r *Request) Delete() *Response {
	return r.Do(DeleteMethod)
}

// Options is a helper method for Do().
func (r *Request) Options() *Response {
	return r.Do(OptionsMethod)
}

// SetBody encodes and sets the proper headers for the request body from the
// given resource.  The resource is encoded in-memory, so be careful about
// passing a massive object.  You can set the ContentLength and Body properties
// manually.
func (r *Request) SetBody(mtype *mediatype.MediaType, resource interface{}) error {
	r.MediaType = mtype
	r.Header.Set(ctypeHeader, mtype.String())

	if resource == nil {
		return nil
	}

	buf, err := mtype.Encode(resource)
	if err != nil {
		return err
	}

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
