package sawyer

import (
	"github.com/bmizerany/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	})

	client := setup.Client
	user := &TestUser{}
	apierr := &TestError{}

	res := client.Get(user, apierr, "404")
	if res.IsError() {
		t.Fatalf("response errored: %s", res.Error())
	}

	assert.Equal(t, 404, res.StatusCode)
}

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

type TestUser struct {
	Id    uint   `json:"id"`
	Login string `json:"login"`
}

type TestError struct {
	Message string `json:"message"`
}

type SetupServer struct {
	Client *Client
	Server *httptest.Server
	Mux    *http.ServeMux
}

func Setup(t *testing.T) *SetupServer {
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	client, err := NewFromString(srv.URL, nil)
	if err != nil {
		t.Fatalf("Unable to parse %s: %s", srv.URL, err.Error())
	}
	return &SetupServer{client, srv, mux}
}

func (s *SetupServer) Teardown() {
	s.Server.Close()
}
