package sawyer

import (
	"errors"
	"github.com/bmizerany/assert"
	"net/http"
	"testing"
)

func TestSuccessfulResponse(t *testing.T) {
	r := &Response{Response: &http.Response{StatusCode: 200}}
	assert.Equal(t, false, r.AnyError())
	assert.Equal(t, false, r.IsError())
	assert.Equal(t, false, r.IsApiError())
	assert.Equal(t, "", r.Error())
	assert.Equal(t, 200, r.StatusCode)
}

func TestErroringResponse(t *testing.T) {
	r := ResponseError(errors.New("wat"))
	assert.Equal(t, true, r.AnyError())
	assert.Equal(t, true, r.IsError())
	assert.Equal(t, false, r.IsApiError())
	assert.Equal(t, "wat", r.Error())
	assert.Equal(t, 0, r.StatusCode)
}

func TestApiErroringResponse(t *testing.T) {
	r := &Response{
		Response: &http.Response{StatusCode: 404},
		// this is typically set by Do() on a *Request
		isApiError: UseApiError(404),
	}
	assert.Equal(t, true, r.AnyError())
	assert.Equal(t, false, r.IsError())
	assert.Equal(t, true, r.IsApiError())
	assert.Equal(t, "", r.Error())
	assert.Equal(t, 404, r.StatusCode)
}
