package sawyer

import (
	"net/url"
	"testing"
)

var endpoints = map[string]map[string]string{
	"http://api.github.com": map[string]string{
		"user": "http://api.github.com/user",
	},
	"http://api.github.com/api/v1": map[string]string{
		"user": "http://api.github.com/api/v1/user",
	},
}

func TestResolve(t *testing.T) {
	for endpoint, tests := range endpoints {
		client, err := NewFromString(endpoint, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		for relative, result := range tests {
			u, err := url.Parse(relative)
			if err != nil {
				t.Errorf(err.Error())
				break
			}

			abs := client.ResolveReference(u)
			if absurl := abs.String(); result != absurl {
				t.Errorf("Bad absolute URL %s for %s + %s == %s", absurl, endpoint, relative, result)
			}
		}
	}
}
