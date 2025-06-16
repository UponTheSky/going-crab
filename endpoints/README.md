# Register Endpoints
In the last chapter, we registered the endpoint of the root path(`/`) and make a response of string `"Hello, Server!"`. In this chapter, we will cover several patterns of registering API endpoints to a running server.

Just like the previous chapter, run the example code using `go run .` in the current directory if you want. Try hitting the following endpoints using tools like `curl` and see what results come in the responses:
- `/`
- `/hello/go`
- `/hello/rust`
- `/hello/swift`

Now, let's look into the code `server.go`.

## The entire structure of handling API requests
There is a specific structure for handling multiple requests for various endpoints in the `net/http` package. We would like to break down the whole structure into several components in details.

### 1. Handler interface and ServeHTTP 
In `net/http`, every type that handles incoming requests implements [`Handler`](https://pkg.go.dev/net/http#Handler) interface. 

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

For example, the type [`HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc) is [implemented as follows](https://cs.opensource.google/go/go/+/refs/tags/go1.24.4:src/net/http/server.go;l=2293):

```go
type HandlerFunc func(ResponseWriter, *Request)

// here the type HandlerFunc implements ServeHTTP, hence it follows the Handler interface
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

func exampleHandlerFunc(w ResponseWriter, r *Request) {
    // do sth
}

func main() {
    f := HandlerFunc(exampleHandlerFunc) // f is now of the HandlerFunc type, i.e. the Handler type
}
```

So if you see any type that implements this ServeHTTP, that means this type can handle incoming HTTP requests.

### 2.ServeMux - Request Multiplexer
But a request coming from the outside doesn't come to any handler directly. A [multiplexer](https://en.wikipedia.org/wiki/Multiplexer) can be thought as a funnel that intakes inputs from various sources, but streams it to a single destination. In `net/http`, we have [`ServeMux`](https://pkg.go.dev/net/http#ServeMux) struct that follows this concept; it receives multiple requests(inputs) for various endpoints, and hands them over to the internal pattern matcher(destination) so that corresponding appropriate logical processes can be executed.

In the previous chapter, we saw `http.DefaultServeMux`. As the name suggests, it is provided as the default ServeMux, and it allows us to simply invoke package-level functions such as `http.Handle` or `http.HandleFunc`. 

```go
// package-level function, http.HandleFunc
// this function registers the function given as the second argument
// as the Handler type, to http.DefaultServeMux
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, Server!")
})

// here we deliberatetly used `http.DefaultServeMux`, but it could be 
// nil for this case
if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
    fmt.Println(err)
}
```

However, from this chapter we will generate a `ServeMux` and pass it to the server-running function(will be coming up in the next chapter). Note that `ServeMux` also implements `ServeHTTP` function, such that it is also of the `Handler` type.

```go
// [...]

// this function creates a ServeMux as a pointer value
mux := http.NewServeMux()

// [...]

if err := http.ListenAndServe(":8080", mux); err != nil {
    log.Fatal(err)
}
```

### 3. Register `Handler`s to `ServeMux` - Handle, HandleFunc
In order for the `ServeMux` to stream the incoming requests to our `Handler`s, we need to register those `Handler`s. In `ServeMux`, there is a function [`Handle`](https://pkg.go.dev/net/http#ServeMux.Handle) that registers a `Handler` that matches the given string pattern. 

```go
// [...]
mux := http.NewServeMux()

// using HandleFunc
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, Server!")
})

// using Handle, by separately defining a struct implementing the Handler type
golang := &LanguageHandler{language: "go"}
rust := &LanguageHandler{language: "rust"}
swift := &LanguageHandler{language: "swift"}

mux.Handle("/hello/go/{$}", golang)
mux.Handle("/hello/rust/{$}", rust)
mux.Handle("/hello/swift/{$}", swift)

// [...]
```

Here note that we use `ServeMux.HandleFunc`, but this is simply an adapter that lets us to pass a function of signature `func(w http.ResponseWriter, r *http.Request)`. Internally it is type-casted to `HandlerFunc`, which is, of course, the `Handler` type.  

**REMARK**
It is sometimes quite confused with these namings - `Handler`, `HandlerFunc`, `Handle`, and `HandleFunc`. How I remember for distingushing them is as follows:
- `Handle` is a verb - so it "Handle"s `Handler`s, which are nouns
- The suffix `Func` means that it is an adapter of its appending type
    - the method `HandleFunc` adapts `Handle` function
    - the type `HandlerFunc` adapts `Handler` type

## Conclusion
In this chapter, we covered the structure inside `net/http` that handles incoming requests. Sending a request hits a `ServeMux` and the endpoint of the request matches registered `Handler`s that implement `ServeHTTP`, using the method `Handle`.

## Exercise
In fact, the type of the second argument of the function [`http.ListenAndServe`](https://pkg.go.dev/net/http#ListenAndServe) is the `Handler`. 

Try passing to it any `Handler`. What happenes when you make requests with different endpoints?
