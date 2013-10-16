# Sawyer

Sawyer is an HTTP user agent for REST APIs.  It is a spiritual compliment to
the [Ruby sawyer gem](https://github.com/lostisland/sawyer).

![](http://techno-weenie.net/sawyer/images/sawyer.jpeg)

Use this to build clients for HTTP/JSON APIs that behave like the GitHub API.


## Pseudo Usage

Just some ideas, there's no chance the final result will look like this :)

```go
type User struct {
  Login string `json:"login"`
}

class ApiError struct {
  Message strign `json:"message"`
}

client := sawyer.NewFromString("https://api.github.com")

// the GitHub API prefers a vendor media type
client.Headers.Set("Accept", "application/vnd.github+json")

apierr := &ApiError{} // decoded from response body on non-20x responses
user := &User{}
req := client.NewRequest("user/21", apierr)
res := req.Get(user)

// get the user's repositories
apierr := &ApiError{}
repos := new([]Repository)
req := client.NewRequest(res.Relations["repos"], apierr)
res := req.Get(repos)

// post a new user
mtype := mediatype.Parse("application/vnd.github+json")
apierr := &ApiError{}
userInput := &User{Login: "bob"}
userOutput := &User{}
req := client.NewRequest("users", apierr)
req.Body = mtype.Encode(userInput)
res := req.Post(userOutput)
```
