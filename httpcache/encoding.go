package httpcache

import (
	"encoding/gob"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/mediaheader"
	"github.com/lostisland/go-sawyer/mediatype"
	"io"
	"net/http"
)

type Adapter interface {
	Get(string, interface{}) *sawyer.Response
	Set(string, *sawyer.Response, interface{}) error
}

func DecodeFrom(v interface{}, resReader io.Reader, bodyReader io.Reader) *sawyer.Response {
	res := Decode(resReader)

	if v != nil && res.ContentLength > 0 {
		err := res.MediaType.Decode(v, bodyReader)
		if err != nil {
			return sawyer.ResponseError(err)
		}
	}

	return res
}

func EncodeTo(v interface{}, res *sawyer.Response, resWriter io.Writer, bodyWriter io.Writer) error {
	if v != nil && res.ContentLength > 0 {
		reader := io.TeeReader(res.Body, bodyWriter)
		dec, err := res.MediaType.Decoder(reader)
		if err != nil {
			return err
		}

		err = dec.Decode(v)
		if err != nil {
			return err
		}
	}

	err := Encode(res, resWriter)
	if err != nil {
		return err
	}

	return nil
}

// Response is an http.Response that can be encoded and decoded safely.
type response struct {
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

	resCopy := response{
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
	var resCopy *response
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
