# Hello, Server!

Ok, let's start our journey to writing a server in Go, with lot of fun!
We'll begin with the most common and widely used type of server, a HTTP server.
The package [`net/http`](https://pkg.go.dev/net/http) fully supports writing HTTP(s) clients and servers. 

You can run the server by running this command on terminal: `go run .`. Check that the server runs correctly at `localhost:8080` using either your browser or some tools like `curl`.

After checking out that you get the response from the server like below, let's break down the code example in `server.go` together:

```sh
Hello, Server!%                                                    
```

## 1. Register endpoints
Like any other server-side programming, writing a server with `net/http` requires you to register endpoints that clients hit with their HTTP requests. 

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, Server!")
})
```

Here we register a function that is defined on the fly, using the method `HandleFunc`. If you look into the package, there are a family of similar names - `HandlerFunc`, `Handle`, `Handler`, etc. We'll cover those in a later chapter, but for now, let's simply take it for granted, as we have justed started our journey.

### Write your response to Writer
Inside the function, we write the string `"Hello, Server!"` to [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter) interface. This is an extension of the [`Writer`](https://pkg.go.dev/io#Writer) interface defined in the [`io`](https://pkg.go.dev/io) package. Like `HandleFunc` above, we will cover various `Reader`s and `Writer`s later, but here basically we only need to recognize that this is a stream through which our written data goes. 

## 2. Run the server
In `net/http`, running server can be thought of these two separate processes:

1. Listen to the incoming requests.

2. Serve the newly generated connections.

Again, we won't cover these topics right now in this introductory chapter, and for simplicity, the code uses `http.ListenAndServe` function. Note that the first argument is the address of the server. But what is `http.DefaultServeMux` in the second parameter? We'll also cover this part later.

## Summary
So much abstractions for running a server! Of course, running a HTTP server is a very complicated process and requires a lot of setups. Thanks to `net/http`, we could write such server with only four lines of code. In the next chapters, we will look into those abstracted parts one by one and re-write our server in a more sophisticated way. 

## Exercise
Add more endpoints and write any messages you like!
You can consult the pattern [matching rules](https://pkg.go.dev/net/http#hdr-Patterns-ServeMux) in `net/http` and make your messages dynamic, using parts of your endpoints.
