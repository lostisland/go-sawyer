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

func Encode(res *sawyer.Response, writer io.Writer) error {
	enc := gob.NewEncoder(writer)

	resCopy := response{
		Expires:          expiration(res),
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

	return enc.Encode(&resCopy)
}

func Decode(reader io.Reader) *sawyer.Response {
	dec := gob.NewDecoder(reader)
	resCopy := &response{}
	err := dec.Decode(&resCopy)
	httpres := &http.Response{}
	res := &sawyer.Response{BodyClosed: false, Response: httpres}

	if err != nil {
		res.ResponseError = err
		res.BodyClosed = true
		return res
	}

	httpres.Status = resCopy.Status
	httpres.StatusCode = resCopy.StatusCode
	httpres.Proto = resCopy.Proto
	httpres.ProtoMajor = resCopy.ProtoMajor
	httpres.ProtoMinor = resCopy.ProtoMinor
	httpres.Header = resCopy.Header
	httpres.ContentLength = resCopy.ContentLength
	httpres.TransferEncoding = resCopy.TransferEncoding
	httpres.Trailer = resCopy.Trailer
	res.MediaType = &resCopy.MediaType
	return res
}

var DefaultExpirationDuration = time.Hour

// response is an http.Response that can be encoded and decoded safely.
type response struct {
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

func expiration(res *sawyer.Response) time.Time {
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
