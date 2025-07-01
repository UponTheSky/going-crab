# Basic Response
From this chapter and following several chapters, we will focus on writing endpoints that respond to various HTTP requests - mainly `GET`, `POST`, `PATCH`, and `DELETE`. But for handling these requests, we need to know how to write responses and send them back to the clients who made the requests.

This chapter is mainly about writing simple and basic responses. Of course there must be various types of responses and we can't cover all of them. Hence we'll cover writing responses of the followings(with the corresponding endpoints in the example server):
- a simple string: `/string`
- a JSON: `/json`
- a file(such as `.html`): `/html`

Like the previous chapters, run the server by `go run .` and test the endpoints above.

## Writing a response inside `Handler`
Let's check out the [`Handler`](https://pkg.go.dev/net/http#Handler) interface again:

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

So if a server wants to communicate with outside requests, all the interactions occur inside this `ServeHTTP` function, more specifically with `ResponseWriter` and `*Request` arguments. As the reader already have noticed, `ResponseWriter` is for sending a response to a client, whereas `*Request` is for getting information about the request from a client. 

Therefore, in this chapter we will mostly deal with `ResponseWriter` interface type for making responses.

## Writing a simple string
Let's write a very simple response first, a string. In the previous chapters, we used `fmt.Fprint`-like functions with the [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter) as its first argument, inside `Handler`s:

```go
func stringHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "This is a string response")
}
```

This is a function that we use for convenience. But for the learning purpose, let's use `Write` function of `ResponseWriter` directly.

```go
func stringHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is a string response\n"))
}
```

And if we use tools like `curl localhost:8080 -v`, then the response is as expected: 

```sh
# [...]
< HTTP/1.1 200 OK
< Date: Thu, 19 Jun 2025 15:37:28 GMT
< Content-Length: 26
< Content-Type: text/plain; charset=utf-8
< 
This is a string response
# [...]
```

## Writing a JSON
But in many web services data is not a mere plain text. Rather, what we usually see is a huge volume of JSON data with tons of properties from API servers. 

Writing JSON is not as simple, direct, and intuitive as writing a string. Let's break down the process as follows:

1. First, we define a JSON schema for the response we want to make, using Go's `struct`.   
    - Technically, any primitive values such as `int` and `string` can be converted into JSON values but in many cases we send compositive JSON values.
    - Another way of creating a JSON schema is using a `map` object. However, for clients to expect consistent response patterns, it is more appropriate to use a fixed schema using `struct`. 

2. After going through several logical processes, we produce an object of the defined `struct`. 

3. Finally, encode the produced `struct` into the corresponding JSON string.
    - In Go, the word *encode* is expressed as *marshall*. As a corollary, *decode* corresponds to *unmarshall*.

### 1. Define JSON schema using `struct`.
Suppose we send a JSON response like:

```ts
{
    "name": "Leslie Chow"
    "actor": "Ken Jeong",
    "series": [1, 2, 3]
}
```

Then we define its corresponding Go `struct` like:

```go
type Character struct {
	Name   string `json:"name"`
	Actor  string `json:"actor"`
	Series []int  `json:"series"`
}
```

It is pretty similar to how we define a JSON-like object in popular languages like JavaScript or Python. However, there are two things that I want you to notice here:
- The library we'll going to use is [`encoding/json`](https://pkg.go.dev/encoding/json) and it ignores non-exposed fields; i.e. we have to specify the fields with their first letters as capital letters. 
- See the raw string literals(surronded by two backticks) right next to the type of each of the field. It tells the library function that its corresponding JSON value has the fields with the names specified as the one in the raw string literals. So here in this case, the encoded JSON will have fields `"name"`, `"actor"`, and `"series"`.

Once you define a schema through `struct`, it's time to encode it into a JSON string.

### 2. Encode `struct` into JSON string
In the `encoding/json` package, there are two ways to encode a `struct` into a JSON string.

1. The most straightforward method is to use [`json.Marshall`](https://pkg.go.dev/encoding/json#Marshal):

```go
chow := Character{
    Name:   "Leslie Chow",
    Actor:  "Ken Jeong",
    Series: []int{1, 2, 3},
}

chow_json, err := json.Marshal(chow)

if err != nil {
    log.Fatal(err)
}
```

Here `chow_json` is an encoded array of byte whose string representation is `{"name":"Leslie Chow","actor":"Ken Jeong","series":[1,2,3]}`. Now we only have to write it to `ResponseWriter` for sending the data to clients:

```go
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	chow := Character{
		Name:   "Leslie Chow",
		Actor:  "Ken Jeong",
		Series: []int{1, 2, 3},
	}

	chow_json, err := json.Marshal(chow)

	if err != nil {
		log.Fatal(err)
	}

	w.Write(chow_json)
}
```

Check the result using tools like `curl`, against the API endpoint `/json`. You will get exactly the same string response as the JSON encoded byte array represents. 

2. However, there is a more convenient and more configurable way of making a JSON response. When a writing stream to which strings would be written(in our case, it is `ResponseWriter`), use [`Encoder.Encode`](https://pkg.go.dev/encoding/json#Encoder.Encode):

```go
// wrap the writing stream as Encoder
encoder := json.NewEncoder(w)

// invoke Encode and directly pass the struct object to the function
if err := encoder.Encode(chow); err != nil {
    log.Fatal(err)
}
```

I personally think the second method is more succinct, easy to write, and more configurable(please check out the `encoding/json` docs). However, encoding to a JSON string may be necessary in some cases, so we need to figure out the problem and choose the right method for it.

## Writing a file
Serving a static file in the server's directory can be handled by `net/http`'s [`FileServer`](https://pkg.go.dev/net/http#FileServer):

```go
// [...]
// note that in the same directory as `server.go`, we have a directory "static"
// that is why we here specify the directory for the static files as `.static`
staticFileHandler := http.FileServer(http.Dir("./static"))
mux.Handle("/", staticHandler)
// [...]
```

Note that we use `ServerMux.Handle` instead of `ServeMux.HandleFunc`, since `htmlHandler` is already a `Handler`. Now if we use `curl` or type in the address `localhost:8080/hangover.html`, you will see the html response either on your terminal or browser.

**REMARK**: `FileServer` vs [`ServeFile`](https://pkg.go.dev/net/http#ServeFile)
There is another similar function called `http.ServeFile` in `net/http` package. Under the hood, they share the same internal logic, but 
- `FileServer` goes through a few redirect check processes, and we don't have to provide the entire file/directory paths of the files we want to serve
- `ServeFile` doesn't go through the same redirect check processes that `FileServer` does, and the user has the responsibility of specifying the full file path, which might have a potential security issue. 

You may find [this old Stack Overflow page](https://stackoverflow.com/questions/28793619/golang-what-to-use-http-servefile-or-http-fileserver) for comparing these two. However, note that the page is old and Go standard libraries have changed a lot since then. 

## Conclusion
We have covered basic methods of writing HTTP responses which are very common in today's Web. Now at least you can send clients any messages you want to say! 

## Exercise
1. Define a bit more complicated `struct`. Probably nested one like this:

```go
type Character struct {
	Name Name
	Actor Actor
	Series: []Series
}

type Name struct {
	First string
	Last string
}

type Actor struct {
	Name Name
	Age int
	Nationality string
}

type Series struct {
	Name string
	ReleaseYear int
}
```

Create an endpoint serving nested JSON responses.

2. Create an endpoint that serves a file other than `.html`, such as `.csv`, `.pdf`, or `.txt`. See how each of them is rendered on your browser.
