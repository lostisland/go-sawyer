package mediatype

import (
	"fmt"
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

func (m *MediaType) Decode(v interface{}, body io.Reader) error {
	if body == nil {
		return nil
	}

	dec := m.Decoder(body)
	if dec == nil {
		return fmt.Errorf("No decoder found for format %s (%s)", m.Format, m.String())
	}

	return dec.Decode(v)
}
