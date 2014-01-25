/*
Package sawyer is a net/http hypermedia client.  It is specifically designed to
talk to APIs designed like the GitHub API.  It includes media type based content
parsing, flexible hypermedia handling, and caching.  It is meant to be used in
an API client, such as Octokit.

The first principle of sawyer is that the sawyer.Client contains common HTTP
properties that are applied to all requests: a base URI, headers, URI query
values, and the Cacher.

Initialize a sawyer.Client with the base URI:

    // the second argument is an optional custom http.Client.
    cli, err := sawyer.NewFromString("https://api.github.com", nil)
    // err would from url.Parse()
    cli.Header.Set("Accept", "application/vnd.github.v3+json")

This client can make an API call like so:

    req, err := cli.NewRequest("users/lostisland")
    // err can be from url.Parse() or http.NewRequest()

    res := req.Get()

The Response type doubles as a potential Go error, or an API error.  An API
error is defined as any response that includes a status that isn't in the 2xx
series.  Here's how one might handle the various response states:

    res := req.Get()
    if res.IsError() {
      // Was there a Go error processing the HTTP request?
      // res.Error() works as expected.
      panic(res)

    } else if res.IsApiError() {
      // Was there an unexpected status code, like 404?
      panic(res.Status)

    } else {
      // All good, cap'n!
    }

Once a valid response has been returned, resources can be deserialized.  The
decoder is set from the MediaType, which is parsed from the response's
Content-Type header by the mediatype package.

    user := &User{}
    err := res.Decode(user)
    fmt.Println(user.Login)

*/
package sawyer
