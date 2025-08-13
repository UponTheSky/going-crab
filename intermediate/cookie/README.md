# Cookie
A HTTP cookie is a piece of data that a HTTP server sends to a browser, such that the browser can use it for various purposes. Therefore, the concept of cookie innately involves two distinct entities - browser and server.

- Server: sends pieces of text data through the `Set-Cookie` header, with several optional flags.
- Browser: sees the cookie(s) in the response from the server, and store them. It can resend the cookie data via the `Cookie` header. 

Note that a HTTP communication between two HTTP servers doesn't require any cookies. You can simply parse body data. However, when your server interacts with browser, you would need to know how to send cookies. 

For detailed usage of cookies, please refer to MDN's [guide page](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Cookies). Here in this chapter we would talk about how to send cookies using `net/http` package.

## Send Cookies over HTTP Response
Although we could set `Set-Cookie` header manually, it is very error prone, because there are a set of rules that you need to follow for setting cookies. Therefore, `net/http` package provides [`Cookie`](https://pkg.go.dev/net/http#Cookie) type and [SetCookie](https://pkg.go.dev/net/http#SetCookie) function for easy cookie handling. Notice that there are more types and functions for cookies in `net/http`, but those are for client-side requests, thus we don't cover them in this book.

First, let's setup a very simple server. We won't separate `Listen` and `Serve` this time, and use the default `ServeMux` for simplicity.

```go
func main() {
	http.HandleFunc("/cookie", func(w http.ResponseWriter, r *http.Request) {
		// set cookie here

		fmt.Fprintln(w, "cookies are set!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Next, we create and set a few cookies inside `HandleFunc`.

```go
// [...]
http.HandleFunc("/cookie", func(w http.ResponseWriter, r *http.Request) {
    // set cookie here
    kenCookie := &http.Cookie{
        Name:   "hangover",
        Value:  "Ken Jeong",
    }

    zachCookie := &http.Cookie{
        Name:   "hangover",
        Value:  "Zach Galifianakis",

    }

    http.SetCookie(w, kenCookie)
    http.SetCookie(w, zachCookie)

    fmt.Fprintln(w, "cookies are set!")
})

// [...]
```

Note that there are multiple options in the [`Cookie`](https://pkg.go.dev/net/http#Cookie) type that define behaviors of a cookie, but we won't cover all of them. Here we set the name and value of our cookie. A HTTP cookie is actually a key-value pair with bunch of options. But it also allows for one key to have multiple values. Hence the above example sets two different cookie values for the same key `"hangover"`. 

Now run the server and try testing with `curl localhost:8080/cookie -v`  on your terminal, and then the response will be like:

```sh
< HTTP/1.1 200 OK
< Set-Cookie: hangover="Ken Jeong"
< Set-Cookie: hangover="Zach Galifianakis"
< Date: Wed, 13 Aug 2025 07:49:23 GMT
< Content-Length: 17
< Content-Type: text/plain; charset=utf-8
< 
cookies are set!
```

## Conclusion
Setting a cookie using `net/http` is very straightforward. You may find it useful in handling user session and other usages. 

## Exercise
Well, this chapter is pretty short, and I've got no exercises to give. Maybe in the project session?
