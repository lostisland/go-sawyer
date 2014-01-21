package sawyer

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer/hypermedia"
	"github.com/lostisland/go-sawyer/mediatype"
	"net/http"
	"strings"
	"testing"
)

// see sawyer_test.go for definitions of structs and SetupServer

func TestSuccessfulGet(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		head := w.Header()
		head.Set("Content-Type", "application/json")
		link := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next", <https://api.github.com/user/repos?page=50&per_page=100>; rel="last"`
		head.Set("Link", link)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"login": "sawyer",
			"url": "/field/self",
			"foo_url": "/field/foo",
			"whatever": "/field/whatevs",
			"homepage_url": "/not/hypermedia",
			"_links": {
				"self": { "href": "/hal/self" },
				"foo": { "href": "/hal/foo" },
				"boom": { "href": "/hal/boom" }
			}
		}`))
	})

	client := setup.Client
	user := &TestUser{}

	req, err := client.NewRequest("user")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, false, res.IsApiError())

	rels := hypermedia.Rels(res)
	assert.Equal(t, 2, len(rels))
	assert.Equal(t, "https://api.github.com/user/repos?page=3&per_page=100", string(rels["next"]))
	assert.Equal(t, "https://api.github.com/user/repos?page=50&per_page=100", string(rels["last"]))

	assert.Equal(t, nil, res.Decode(user))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "sawyer", user.Login)

	userRels := hypermedia.Rels(user)
	assert.Equal(t, 5, len(userRels))
	assert.Equal(t, "/hal/self", string(userRels["self"]))
	assert.Equal(t, "/hal/foo", string(userRels["foo"]))
	assert.Equal(t, "/hal/boom", string(userRels["boom"]))
	assert.Equal(t, "/field/whatevs", string(userRels["whatevs"]))
	assert.Equal(t, "/field/self", string(userRels["Url"]))
}

func TestSuccessfulGetWithoutOutput(t *testing.T) {
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
	user := &TestUser{HALResource: &hypermedia.HALResource{}}

	req, err := client.NewRequest("user")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, false, res.IsApiError())

	assert.Tf(t, !res.IsError(), "Response shouldn't have error")
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, false, res.BodyClosed)
	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "", user.Login)

	dec := json.NewDecoder(res.Body)
	dec.Decode(user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "sawyer", user.Login)
}

func TestSuccessfulGetWithoutDecoder(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		head := w.Header()
		head.Set("Content-Type", "application/booya+booya")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "login": "sawyer"}`))
	})

	client := setup.Client
	user := &TestUser{HALResource: &hypermedia.HALResource{}}

	req, err := client.NewRequest("user")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.NotEqual(t, nil, res.Decode(user), "response should have decoder error")
	assert.Tf(t, strings.HasPrefix(res.Error(), "No decoder found for format booya"), "Bad error: %s", res.Error())
}

func TestSuccessfulPost(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	mtype, err := mediatype.Parse("application/json")

	setup.Mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, mtype.String(), r.Header.Get("Content-Type"))

		user := &TestUser{HALResource: &hypermedia.HALResource{}}
		mtype.Decode(user, r.Body)
		assert.Equal(t, "sawyer", user.Login)

		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"login": "sawyer2"}`))
	})

	client := setup.Client
	user := &TestUser{HALResource: &hypermedia.HALResource{}}

	req, err := client.NewRequest("users")
	assert.Equal(t, nil, err)

	user.Login = "sawyer"
	req.SetBody(mtype, user)
	res := req.Post()
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, false, res.IsApiError())
	assert.Equal(t, nil, res.Decode(user))

	assert.Equal(t, nil, err)
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, "sawyer2", user.Login)
	assert.Equal(t, true, res.BodyClosed)
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
	user := &TestUser{HALResource: &hypermedia.HALResource{}}
	apierr := &TestError{}

	req, err := client.NewRequest("404")
	if err != nil {
		t.Fatalf("request errored: %s", err)
	}

	res := req.Get()
	assert.Equal(t, true, res.IsApiError())
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, nil, res.Decode(apierr))

	assert.Equal(t, 404, res.StatusCode)
	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "", user.Login)
	assert.Equal(t, "not found", apierr.Message)
	assert.Equal(t, true, res.BodyClosed)
}

func TestResolveRequestQuery(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/q", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "1", q.Get("a"))
		assert.Equal(t, "4", q.Get("b"))
		assert.Equal(t, "3", q.Get("c"))
		assert.Equal(t, "2", q.Get("d"))
		assert.Equal(t, "1", q.Get("e"))
		w.WriteHeader(123)
		w.Write([]byte("ok"))
	})

	assert.Equal(t, "1", setup.Client.Query.Get("a"))
	assert.Equal(t, "1", setup.Client.Query.Get("b"))

	setup.Client.Query.Set("b", "2")
	setup.Client.Query.Set("c", "3")

	req, err := setup.Client.NewRequest("/q?d=4")
	assert.Equal(t, nil, err)

	req.Query.Set("b", "4")
	req.Query.Set("c", "3")
	req.Query.Set("d", "2")
	req.Query.Set("e", "1")

	res := req.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, 123, res.StatusCode)
}
