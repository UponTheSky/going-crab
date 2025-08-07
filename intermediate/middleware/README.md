# Middleware
As a server application gets bigger, then there would be several core business logics inside. For instance, when it comes to building a social platform, there must be parts about users, posts, media, analytics, to name a few. And handling each part involves several common neceesary logics such as user auth and logging. 

This is where something called *middleware* comes in. A middleware sits between a user request and a collection of business units of an application, and handles common logics like logging, which must be duplicated for each of the unit if otherwise.

In this chapter, we will look into how to write a middleware in Go, only using `net/http`, which means, this chapter would be pretty short.

## A Middleware is Another `handler`
If you recall the [endpoints](../../basic/endpoints/README.md) chapter, we have covered what is the `handler` type and how we cooperate it with other functions and struct types in `net/http` to write a server.

To tell you from the conclusion, a middleware is another `handler`, where it receives an HTTP request and hands it over to other `handler`s handling actual business problems. Of course, as a middleware, there must be its own logic that is to be executed *before* and/or *after* the other `handler` takes over the request. 

```go
func MyMiddleware(w http.ResponseWriter, r *http.Request) {
    // before: handle your own logic

    someHandler.ServeHTTP(w, r)

    // after: handle your own logic
}
```

Of course, this function needs to be type-changed through `http.HandlerFunc`.

But here are a few problems in the example code above. This `MyMiddleware` handler function needs to have `someHandler` from its same or outer scopes. Hence, if this function is defined at the global scope, `someHandler` also needs to be defined at the global scope. This reduces our flexibility of writing code in places we want.

The other problem is that you can't generalize the handlers inside. What about `otherHandler`? Would you define the middlewares for each and every handlers defined in the application? 

To resolve this issue, we would like to define a factory method instead, where the factory produces a middleware, given a `handler`. 

```go
func AddLoggingMiddleware(handler http.Handler) http.Handler {
	middlware := func(w http.ResponseWriter, r *http.Request) {
		// log the method of the request, path(endpoint), and the current time
		log.Println(r.Method, r.URL.Path, time.Now())

		// hands the request over to the given handler
		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(middlware)
}
```

At first the code looks a little bit complicated, but it is not. `AddLoggingMiddleware` is only making a `handler` given another `handler` as its input argument. Inside this factory function, we define a middleware, which has the exact same form as `MyMiddleware`. Inside this "middleware", we pass the request to the given handler.

That's it! We just defined a new middleware where it makes a log about the HTTP method and the endpoint of the given request, and the current time where the request is passed to. It is very simple but also an elegant way of using the type system of `net/http`.

## Where to Put a Middleware?
Where to put a middleware? It is totally upto the developer. If it is about logging, it is good for it to wrap the uppermost `ServeMux`(recall that a `ServeMux` is also of `handler` type), because it needs to log all the information of the incoming requests. However, you can narrow down the scope. If you have a middleware about a user auth and you also provide a public API that is open to anyone, you would only need to apply that middleware to user-specific features of the application, which are also represented as `handler`s. 

Since we have defined a logging middleware, let's put the middleware to the uppermost `ServeMux`:

```go
func main() {
	// mux and handlers
	mux := http.NewServeMux()

	mux.HandleFunc("/chow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Kaman! Kachick!")
	})

	mux.HandleFunc("/alan", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey Alan, what are you doing up there!")
	})

	// add our middleware here
	mux = AddLoggingMiddleware(mux).(*http.ServeMux)

	// listener
	listener, err := net.Listen("tcp", ":8080") // if you need to pass context, use ListeConfig.Listen

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	// server
	server := &http.Server{Handler: mux}
	defer server.Close()

	err = server.Serve(listener)
	log.Fatal(err)
}
```

Running the server and making http requests to `/chow` or `/alan` will show you log messages on your terminal like:

```sh
2025/08/07 12:30:06 GET /chow 2025-08-07 12:30:06.389967 +0900 JST m=+1.386161418
2025/08/07 12:30:33 GET /alan 2025-08-07 12:30:33.42189 +0900 JST m=+28.418270251
```

## Conclusion
Hurray! Now we know how to write and add our own middlewares. We will write a handful of middlewares in the project section of this chapter soon. So stay tuned!

## Exercise
1. Try make our logging middleware more enriched with information about the incoming request. What more information can we get from [http.Request](https://pkg.go.dev/net/http#Request) type? How about generating our own session ID and pass it to the context of the request(this is a bit advanced topic though! We may cover the topic about session, but has not planned yet).

2. Try write other middlewares. I can't think of one immediately now, so it's upto you the readers, to expand your creativity :)
