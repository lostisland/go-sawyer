package sawyer

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessfulGet(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "login": "sawyer"}`))
	})

	client := setup.Client
	user := &TestUser{}
	apierr := &TestError{}

	req, err := client.NewRequest("user", apierr)
	if err != nil {
		t.Fatalf("request errored: %s", err)
	}

	res, err := req.Get(user)
	if err != nil {
		t.Fatalf("response errored: %s", err)
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "sawyer", user.Login)
	assert.Equal(t, "", apierr.Message)
}

func TestSuccessfulPost(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	mtype, err := mediatype.Parse("application/json")

	setup.Mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, mtype.String(), r.Header.Get("Content-Type"))

		user := &TestUser{}
		mtype.Decode(user, r.Body)
		assert.Equal(t, "sawyer", user.Login)

		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"login": "sawyer2"}`))
	})

	client := setup.Client
	user := &TestUser{}
	apierr := &TestError{}

	req, err := client.NewRequest("users", apierr)
	if err != nil {
		t.Fatalf("request errored: %s", err)
	}

	user.Login = "sawyer"
	req.SetBody(mtype, user)
	res, err := req.Post(user)
	if err != nil {
		t.Fatalf("response errored: %s", err)
	}

	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, "sawyer2", user.Login)
	assert.Equal(t, "", apierr.Message)
}

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

	req, err := client.NewRequest("404", apierr)
	if err != nil {
		t.Fatalf("request errored: %s", err)
	}

	res, err := req.Get(user)
	if err != nil {
		t.Fatalf("response errored: %s", err)
	}

	assert.Equal(t, 404, res.StatusCode)
	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "", user.Login)
	assert.Equal(t, "not found", apierr.Message)
}

type TestUser struct {
	Id    int    `json:"id"`
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
