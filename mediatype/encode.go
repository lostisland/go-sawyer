package mediatype

import (
	"bytes"
	"fmt"
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
	if v == nil {
		return nil, fmt.Errorf("Nothing to encode")
	}

	buf := new(bytes.Buffer)
	enc := m.Encoder(buf)
	if enc == nil {
		return buf, fmt.Errorf("No encoder found for format %s (%s)", m.Format, m.String())
	}

	return buf, enc.Encode(v)
}
