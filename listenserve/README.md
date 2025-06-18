# Listen And Serve

For sure, registering `Handler`s to `ServeMux` doesn't do anything; you need to pass it to functions that run actual server processes so that any external HTTP requests can arrive at the server in real-time.

In this chapter, we will cover how `net/http` can communicate with the outside world so that requests from the world can be read by registered `Handler`s in the last chapter.

There would be nothing new for returning responses from our server in this chapter, but for those who are interested, run the command `go run .` to run the server written in `server.go`.

Note that we won't cover much about computer networking or socket programming in details, albeit some mentions on them may be necessary. The packages `net` or `net/http` doesn't expose low-level socket interface directly, but they wrap it in nice ways. Thus we would like to look over the entire picture rather than each single detailed low-level API.

## 1. Listen - address and protocol

In unix-like socket programming, we need the following two pieces of information to run a server process. First is the **address**(*host* + *port*) of the server that we want the requests to arrive at. Second is the type of the protocol under which we want the network communication to be held.

**REMARK**: In `net/http`, the TCP protocol is set by default in the [`Server.ListenAndServe`](https://pkg.go.dev/net/http#Server.ListenAndServe) function. For our purpose of this chapter(and throughout the entire book), choosing TCP will be sufficient.

In the package `net`, setting the address and the protocol occurs when generating a [`Listener`](https://pkg.go.dev/net#Listener). Note that creating a listener requires the package `net`, not `net/http`. The naming is fairly intuitive; the server **listens** to the outside world to send HTTP requests for the given **address**, following the **protocol**.

```go
// [...]

// Note that these lines of code will not be showing up in server.go 

listener, err := net.Listen("tcp", ":8080") // if you need to pass context, use ListeConfig.Listen

if err != nil { 
    log.Fatal(err) 
}

defer listener.Close() // don't forget to retrieve the resource before the function finishes 
```

Here we use the package-level function [`net.Listen`](https://pkg.go.dev/net#Listen), a nice wrapper of the functions of the various protocols such as [`net.ListenTCP`](https://pkg.go.dev/net#ListenTCP). We could have just used this function for creating a listener specific to the TCP protocol without much benefits. Hence, from this chapter, we'll stick to this nice `net.Listen`.

After creating a listener object, we could directly write a simple server as follows:

```go
// [...]

for {
    conn, err := listener.Accept()

    if err != nil {
        log.Println(err)
        continue
    }

    go func(conn net.Conn) {
        defer conn.Close()

        conn.Write([]byte("Hey client there!"))
    }(conn)
}
```

When you use tools like `nc` in case of macOS, this will make the response `"Hey client there!"`(note that tools like curl would not work, since this server doesn't follow the HTTP protocol after HTTP1.0). However, this will be a bit different topic from our journey, more relatable in computer networking classes. Fortunately, `net/http` provides a convenient struct [`net.Server`](https://pkg.go.dev/net/http#Server) to handle this infinite loop of handling request for you.

## 2. Serve

By using the package `net/http`, creating a new server instance is easy. But don't forget to prepare a `ServeMux` first!

```go
// [...]

server := &http.Server{Handler: mux}
defer server.Close()
```

Before "serving" requests from outside, `Server` requires a `Handler` to be registered beforehand. For handling multiple endpoints, we need a `ServeMux` in the place for the Handler.

By the way, there are multiple options configurable inside `http.Server`. We will cover those in later chapters.

After the server is ready, now it's time to run the server.

```go
// [...]
if err := server.Serve(listener); err != nil {
    log.Fatal(err)
}
```

That...'s it! This is all. How simple and beautiful the `net/http` package make abstractions on the low-level socket programming!

## Conclusion

But to this far I guess the readers may already have noticed; This whole process could have been simply done by one single method, `net.ListenAndServe` that has showed up in the past two chapters. This is true, because if you see the source code of this function, it goes through exactly the same process that we have followed so far. Then why have we separately written down the parts in details?

1. First, we separate `Server`. By doing so, we can put it under our control with more specific configurations on our server.

2. Second, we separate `Listener`. We don't have to be confined to using the TCP protocol but other options are also available such as UDP.

From the last chapter and the current chapter, we have looked at how a simple server works and its structure using the package `net/http`. From the next chapter, we would like to try writing actual HTTP responses with some tests.

## Exercise

No exercise today. You've done a great job today by reading this book to this far. Well done! Take some rest for the next chapters which will be a bit more exciting and fun :smile:.
