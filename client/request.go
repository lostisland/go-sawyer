package sawyer

import (
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Client   *http.Client
	ApiError interface{}
	*http.Request
}

const (
	GetMethod  = "GET"
	PostMethod = "POST"
)

func (c *Client) NewRequest(rawurl string, apierr interface{}) (*Request, error) {
	u, err := c.resolveReferenceString(rawurl)
	if err != nil {
		return nil, err
	}

	httpreq, err := http.NewRequest(GetMethod, u, nil)
	return &Request{c.HttpClient, apierr, httpreq}, err
}

func (r *Request) Get(output interface{}) (*Response, error) {
	return r.Do(GetMethod, output)
}

func (r *Request) Post(output interface{}) (*Response, error) {
	return r.Do(PostMethod, output)
}

func (r *Request) SetBody(mtype *mediatype.MediaType, input interface{}) error {
	buf, err := mtype.Encode(input)
	if err != nil {
		return err
	}
	r.ContentLength = int64(buf.Len())
	r.Body = ioutil.NopCloser(buf)
	return nil
}
