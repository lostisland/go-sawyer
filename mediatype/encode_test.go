package mediatype

import (
	"github.com/bmizerany/assert"
	"io"
	"testing"
)

func TestAddEncoder(t *testing.T) {
	AddEncoder("test", func(w io.Writer) Encoder {
		return &PersonEncoder{w}
	})

	mt, err := Parse("application/test+test")
	if err != nil {
		t.Fatalf("Error parsing media type: %s", err.Error())
	}

	person := &Person{"bob"}
	buf, err := mt.Encode(person)
	if err != nil {
		t.Fatalf("Error encoding: %s", err.Error())
	}

	assert.Equal(t, "bob", buf.String())
}

type PersonEncoder struct {
	body io.Writer
}

func (d *PersonEncoder) Encode(v interface{}) error {
	if p, ok := v.(*Person); ok {
		d.body.Write([]byte(p.Name))
	}
	return nil
}
