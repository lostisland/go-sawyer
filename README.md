# Sawyer

Sawyer is an HTTP user agent for REST APIs.  It is a spiritual compliment to
the [Ruby sawyer gem](https://github.com/lostisland/sawyer).

![](http://techno-weenie.net/sawyer/images/sawyer.jpeg)

Use this to build clients for HTTP/JSON APIs that behave like the GitHub API.


## Pseudo Usage

```go
type User struct {
  Login string `json:"login"`
}

class ApiError struct {
  Message strign `json:"message"`
}

client := sawyer.NewFromString("https://api.github.com")

// this sets Accept header to application/json
client.SetEncoding("json")

// the GitHub API prefers a vendor media type
client.Headers.Set("Accept", "application/vnd.github+json")

// this is the struct that decodes JSON errors
client.SetError(ApiError)

user := &User{}
res := client.Get(user, "user/21")

// get the user's repositories
repos := []Repository{}
res2 := client.Get(repos, res.relations["repos"])
```
