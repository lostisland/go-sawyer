package httpcache

import (
	"encoding/gob"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/mediaheader"
	"github.com/lostisland/go-sawyer/mediatype"
	"io"
	"net/http"
)

// Response is an http.Response that can be encoded and decoded safely.
type Response struct {
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

func Encode(res *sawyer.Response, writer io.Writer) error {
	enc := gob.NewEncoder(writer)

	resCopy := Response{res.Status, res.StatusCode, res.Proto, res.ProtoMajor,
		res.ProtoMinor, res.Header, res.ContentLength, res.TransferEncoding,
		res.Trailer, *res.MediaType}

	return enc.Encode(&resCopy)
}

func Decode(reader io.Reader) *sawyer.Response {
	dec := gob.NewDecoder(reader)
	var resCopy *Response
	err := dec.Decode(&resCopy)
	if err != nil {
		return sawyer.ResponseError(err)
	}

	httpres := http.Response{
		Status:           resCopy.Status,
		StatusCode:       resCopy.StatusCode,
		Proto:            resCopy.Proto,
		ProtoMajor:       resCopy.ProtoMajor,
		ProtoMinor:       resCopy.ProtoMinor,
		Header:           resCopy.Header,
		ContentLength:    resCopy.ContentLength,
		TransferEncoding: resCopy.TransferEncoding,
		Trailer:          resCopy.Trailer,
	}

	var headerDecoder mediaheader.Decoder
	return &sawyer.Response{
		MediaType:   &resCopy.MediaType,
		MediaHeader: headerDecoder.Decode(httpres.Header),
		BodyClosed:  false,
		Response:    &httpres,
	}
}
