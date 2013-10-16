package mediatype

import (
	"io"
)

var decoders = make(map[string]DecoderFunc)

type DecoderFunc func(r io.Reader) Decoder

type Decoder interface {
	Decode(v interface{}) error
}

func AddDecoder(format string, decfunc DecoderFunc) {
	decoders[format] = decfunc
}

func (m *MediaType) Decoder(body io.Reader) Decoder {
	if decfunc, ok := decoders[m.Format]; ok {
		return decfunc(body)
	}
	return nil
}
