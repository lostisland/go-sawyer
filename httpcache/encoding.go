package httpcache

import (
	"bytes"
	"encoding/gob"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/mediatype"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Encode will create a CachedResponse from the sawyer Response, and encode it
// to the given writer.
func Encode(res *sawyer.Response, writer io.Writer) error {
	resCopy := CachedResponse{
		Expires:          expiration(res.Response),
		Status:           res.Status,
		StatusCode:       res.StatusCode,
		Proto:            res.Proto,
		ProtoMajor:       res.ProtoMajor,
		ProtoMinor:       res.ProtoMinor,
		Header:           res.Header,
		ContentLength:    res.ContentLength,
		TransferEncoding: res.TransferEncoding,
		Trailer:          res.Trailer,
	}

	if res.MediaType != nil {
		resCopy.MediaType = *res.MediaType
	}

	return EncodeResponse(&resCopy, writer)
}

// EncodeResponse encodes the CachedResponse to the given writer.
func EncodeResponse(cached *CachedResponse, writer io.Writer) error {
	return gob.NewEncoder(writer).Encode(cached)
}

// EncodeBody copies the response's Body to the given writer.
func EncodeBody(res *sawyer.Response, bodyWriter io.Writer) error {
	if res.ContentLength == 0 {
		return nil
	}

	buf := &bytes.Buffer{}
	writer := io.MultiWriter(bodyWriter, buf)
	_, err := io.Copy(writer, res.Body)
	if err == nil {
		res.Body = ioutil.NopCloser(buf)
	}

	return err
}

// Decode decodes the CachedResponse from the given reader.  It is then wrapped
// by a CachedResponseDecoder that is able to turn the CachedResponse data
// to a sawyer Response.
func Decode(reader io.Reader) (*CachedResponseDecoder, error) {
	dec := gob.NewDecoder(reader)
	res := &CachedResponse{}
	if err := dec.Decode(&res); err != nil {
		return nil, err
	}

	return &CachedResponseDecoder{CachedResponse: res}, nil
}

var DefaultExpirationDuration = time.Hour

// CachedResponse is an http.Response that can be encoded and decoded safely.
type CachedResponse struct {
	Expires          time.Time
	Status           string // e.g. "200 OK"
	StatusCode       int    // e.g. 200
	Proto            string // e.g. "HTTP/1.0"
	ProtoMajor       int    // e.g. 1
	ProtoMinor       int    // e.g. 0
	Header           http.Header
	ContentLength    int64
	TransferEncoding []string
	Trailer          http.Header
	MediaType        mediatype.MediaType
}

// CachedResponseDecoder can decode the embedded CachedResponse into a sawyer
// response.
type CachedResponseDecoder struct {
	Cacher      sawyer.Cacher
	SetBodyFunc func(res *sawyer.Response)
	*CachedResponse
}

// Decode converts the embedded CachedResponse to a sawyer Response.
func (r *CachedResponseDecoder) Decode(req *sawyer.Request) *sawyer.Response {
	cached := r.CachedResponse
	res := &sawyer.Response{
		BodyClosed: false,
		MediaType:  &cached.MediaType,
		Response: &http.Response{
			Status:           cached.Status,
			StatusCode:       cached.StatusCode,
			Proto:            cached.Proto,
			ProtoMajor:       cached.ProtoMajor,
			ProtoMinor:       cached.ProtoMinor,
			Header:           cached.Header,
			ContentLength:    cached.ContentLength,
			TransferEncoding: cached.TransferEncoding,
			Trailer:          cached.Trailer,
			Request:          req.Request,
		},
	}

	res.Cacher = r.Cacher
	if res.Cacher == nil {
		res.Cacher = req.Cacher
	}

	if r.SetBodyFunc != nil {
		r.SetBodyFunc(res)
	}

	return res
}

// IsExpired returns true if the CachedResponse needs to be refreshed.
func (r *CachedResponseDecoder) IsExpired() bool {
	return time.Now().After(r.Expires)
}

// IsFresh returns true if the CachedResponse does not need to be refreshed.
func (r *CachedResponseDecoder) IsFresh() bool {
	return !r.IsExpired()
}

// SetupRequest passes the cached ETag and Last Modified date to the request.
func (r *CachedResponseDecoder) SetupRequest(req *http.Request) {
	if etag := r.Header.Get(etagHeader); len(etag) > 0 {
		req.Header.Set(ifNoneMatchHeader, etag)
	}

	if lastmod := r.Header.Get(lastModHeader); len(lastmod) > 0 {
		req.Header.Set(ifModSinceHeader, lastmod)
	}
}

func expiration(res *http.Response) time.Time {
	return time.Now().Add(maxAgeDuration(res.Header.Get("Cache-Control")))
}

func maxAgeDuration(header string) time.Duration {
	if len(header) > 0 {
		for _, field := range strings.Fields(header) {
			pieces := strings.SplitN(field, "=", 2)
			if len(pieces) != 2 || pieces[0] != "max-age" {
				continue
			}
			if dur, err := time.ParseDuration(pieces[1] + "s"); err == nil {
				return dur
			}
		}
	}

	return DefaultExpirationDuration
}

const (
	etagHeader        = "ETag"
	lastModHeader     = "Last-Modified"
	ifNoneMatchHeader = "If-None-Match"
	ifModSinceHeader  = "If-Modified-Since"
)
