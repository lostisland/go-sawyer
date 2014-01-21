package sawyer

import (
	"bytes"
	"errors"
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer/hypermedia"
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
	"testing"
)

// see sawyer_test.go for definitions of structs and SetupServer

func TestDecode(t *testing.T) {
	value := decoableValue()
	res, err := decodableResponse("application/json", `{"id":1}`)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, res.BodyClosed)

	err = res.Decode(value)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, value.Id)
	assert.Equal(t, true, res.BodyClosed)
}

func TestDecodeWithBadFormatting(t *testing.T) {
	value := decoableValue()
	res, err := decodableResponse("application/json", `{`)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, res.BodyClosed)

	err = res.Decode(value)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, true, res.BodyClosed)
}

func TestDecodeWithNilValue(t *testing.T) {
	res, err := decodableResponse("application/json", `{"id":1}`)
	assert.Equal(t, nil, err)

	err = res.Decode(nil)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, errorNoResource.Error(), err.Error())
}

func TestDecodeWithResponseError(t *testing.T) {
	msg := "whatevs"
	res, err := decodableResponse("application/json", `{"id":1}`)
	assert.Equal(t, nil, err)
	res.ResponseError = errors.New(msg)

	err = res.Decode(decoableValue())
	assert.NotEqual(t, nil, err)
	assert.Equal(t, msg, err.Error())
}

func TestDecodeWithClosedBody(t *testing.T) {
	res, err := decodableResponse("application/json", `{"id":1}`)
	assert.Equal(t, nil, err)
	res.BodyClosed = true

	err = res.Decode(decoableValue())
	assert.NotEqual(t, nil, err)
	assert.Equal(t, errorBodyClosed.Error(), err.Error())
}

func TestDecodeWithoutMediaType(t *testing.T) {
	value := decoableValue()
	res, err := decodableResponse("text/plain", "")
	res.MediaType = nil
	assert.Equal(t, nil, err)

	err = res.Decode(value)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, errorNoMediaType.Error(), err.Error())
}

func TestDecodeWithBadRelCaching(t *testing.T) {
	value := decoableValue()
	res, err := decodableResponse("application/json", `{"id":1}`)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, res.BodyClosed)

	res.Cacher = &badRelCacher{noOpCache: &noOpCache{}}
	err = res.Decode(value)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "Bad Rels", err.Error())
	assert.Equal(t, nil, res.ResponseError)
	assert.Equal(t, 1, value.Id)
	assert.Equal(t, true, res.BodyClosed)
}

func decoableValue() *TestUser {
	return &TestUser{HALResource: &hypermedia.HALResource{}}
}

func decodableResponse(rawtype, rawbody string) (*Response, error) {
	mt, err := mediatype.Parse(rawtype)
	if err != nil {
		return nil, err
	}

	httpres := &http.Response{
		Body:          ioutil.NopCloser(bytes.NewBufferString(rawbody)),
		ContentLength: int64(len(rawbody)),
	}

	res := NewResponse(httpres)
	res.MediaType = mt

	return res, nil
}

type badRelCacher struct {
	*noOpCache
}

func (c *badRelCacher) SetRels(req *http.Request, rels hypermedia.Relations) error {
	return errors.New("Bad Rels")
}
