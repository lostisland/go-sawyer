package mediatype

import (
	"bytes"
	"github.com/bmizerany/assert"
	"io"
	"io/ioutil"
	"testing"
)

func TestAddDecoder(t *testing.T) {
	AddDecoder("test", func(r io.Reader) Decoder {
		return &PersonDecoder{r}
	})

	buf := bytes.NewBufferString("bob")
	mt, err := Parse("application/test+test")
	if err != nil {
		t.Fatalf("Error parsing media type: %s", err.Error())
	}

	person := &Person{}
	err = mt.Decode(person, buf)
	if err != nil {
		t.Fatalf("Error decoding: %s", err.Error())
	}
	assert.Equal(t, "bob", person.Name)
}

type PersonDecoder struct {
	body io.Reader
}

func (d *PersonDecoder) Decode(v interface{}) error {
	if p, ok := v.(*Person); ok {
		by, err := ioutil.ReadAll(d.body)
		if err != nil {
			return err
		}
		p.Name = string(by)
	}
	return nil
}
