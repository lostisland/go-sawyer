package sawyer

import (
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Client    *http.Client
	ApiError  interface{}
	MediaType *mediatype.MediaType
	*http.Request
}

func (c *Client) NewRequest(rawurl string, apierr interface{}) (*Request, error) {
	u, err := c.ResolveReferenceString(rawurl)
	if err != nil {
		return nil, err
	}

	httpreq, err := http.NewRequest(GetMethod, u, nil)
	return &Request{c.HttpClient, apierr, nil, httpreq}, err
}

func (c *Client) NewRelation(link *Hyperlink, args M, apierr interface{}) (*Request, error) {
	u, err := link.Expand(args)
	if err != nil {
		return nil, err
	}

	return c.NewRequest(u.String(), apierr)
}

func (r *Request) Head(output interface{}) (*Response, error) {
	return r.Do(HeadMethod, output)
}

func (r *Request) Get(output interface{}) (*Response, error) {
	return r.Do(GetMethod, output)
}

func (r *Request) Post(output interface{}) (*Response, error) {
	return r.Do(PostMethod, output)
}

func (r *Request) Put(output interface{}) (*Response, error) {
	return r.Do(PutMethod, output)
}

func (r *Request) Patch(output interface{}) (*Response, error) {
	return r.Do(PatchMethod, output)
}

func (r *Request) Delete(output interface{}) (*Response, error) {
	return r.Do(DeleteMethod, output)
}

func (r *Request) Options(output interface{}) (*Response, error) {
	return r.Do(OptionsMethod, output)
}

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
