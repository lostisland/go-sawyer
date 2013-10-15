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
client.Headers.Set("Accept", "application/vnd.github+json")
client.SetEncoding("json")
client.SetError(ApiError)

user := &User{}
res := client.Get(user, "user/21")
```
