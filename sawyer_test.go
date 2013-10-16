package sawyer

import (
	"github.com/bmizerany/assert"
	"net/url"
	"testing"
)

var endpoints = map[string]map[string]string{
	"http://api.github.com": map[string]string{
		"user":                "http://api.github.com/user",
		"/user":               "http://api.github.com/user",
		"http://api.com/user": "http://api.com/user",
	},
	"http://api.github.com/api/v1": map[string]string{
		"user":                "http://api.github.com/api/v1/user",
		"/user":               "http://api.github.com/user",
		"http://api.com/user": "http://api.com/user",
	},
}

func TestResolve(t *testing.T) {
	for endpoint, tests := range endpoints {
		client, err := NewFromString(endpoint, nil)
		if err != nil {
			t.Fatal(err.Error())
		}

		for relative, result := range tests {
			u, err := url.Parse(relative)
			if err != nil {
				t.Error(err.Error())
				break
			}

			abs := client.ResolveReference(u)
			if absurl := abs.String(); result != absurl {
				t.Errorf("Bad absolute URL %s for %s + %s == %s", absurl, endpoint, relative, result)
			}
		}
	}
}

func TestResolveWithQuery(t *testing.T) {
	client, err := NewFromString("http://api.github.com?a=1&b=1", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "1", client.Query.Get("a"))
	assert.Equal(t, "1", client.Query.Get("b"))

	client.Query.Set("b", "2")
	client.Query.Set("c", "3")
	u, err := client.ResolveReferenceString("/foo?d=4")
	assert.Equal(t, "http://api.github.com/foo?a=1&b=2&c=3&d=4", u)
}
