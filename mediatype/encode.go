package mediatype

import (
	"bytes"
	"io"
)

var encoders = make(map[string]EncoderFunc)

type EncoderFunc func(w io.Writer) Encoder

type Encoder interface {
	Encode(v interface{}) error
}

func AddEncoder(format string, encfunc EncoderFunc) {
	encoders[format] = encfunc
}

func (m *MediaType) Encoder(w io.Writer) Encoder {
	if encfunc, ok := encoders[m.Format]; ok {
		return encfunc(w)
	}
	return nil
}

func (m *MediaType) Encode(v interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if enc := m.Encoder(buf); enc != nil {
		if err := enc.Encode(v); err != nil {
			return buf, err
		}
	}
	return buf, nil
}
