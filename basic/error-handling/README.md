# Error Handling
For sure there must be errors in your application. Logging an error is important for your post-event handling, but you should do something when there is an error. In this chapter, we will briefly cover how to handle errors ocurring in a server using `net/http` package.

In this chapter, check out the file `server.go` and test it with tools like `curl`. 

The endpoints are
- `/error/client/{errorCode}`: returns responses about client-side errors
- `/error/server`: returns responses about server-side errors

## Possible Errors in a Server Application
There could be multiple reasons for errors occurring on server side. It could be a bad request from a client request, a logical error from your own data handling logic, or some low-level error such as a language flaw, a resource usage upsurge, or even a hardware flaw. Obviously, as application engineers, we cannot handle those low-level errors(it's probably DevOps people's). 

- As for errors from users's requests, they can be categorized into several typical cases. A user might ask non-existing information(not found), or she/he has no permission to access certain resources(forbidden). Consult [this MDN page](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status#client_error_responses) for more comprehensive details.

- Server-originated errors are relatively few, but in our case where we write application logic, it would be mostly the internal logical error(internal server error). 

## Handling Errors with `http.Error`
In the package `net/http` we have a very convenient built-in function that returns a HTTP response about an error, [`http.Error`](https://pkg.go.dev/net/http#Error).

The client-side related errors are handled at `/error/client/errorCode` endpoint:

```go
mux.HandleFunc("/error/client/{code}", func(w http.ResponseWriter, r *http.Request) {
    code := r.PathValue("code")

    switch codeParsed {
    case "400":
        http.Error(w, "Bad request", http.StatusBadRequest)
    default:
        http.Error(w, "Not found", http.StatusNotFound)
    }
})
```
When there is a need for the server to respond with information about error, you simply use `http.Error` with a string and the status code provided. The status code is already prepared inside `net/http`, thus you don't have to type in the integer literals everytime.

For sure, this function is opinionated and you can a fully-customized error response depending on your need, as we saw in the previous chapters. 

Very similar to the client-side, the server-side errors are handled as below:

```go
mux.HandleFunc("/error/server", func(w http.ResponseWriter, r *http.Request) {
    err := errors.New("some kind of error you met while executing the internal logic")

    http.Error(w, err.Error(), http.StatusInternalServerError)
})
```

## Conclusion
This chapter is as simple as the previous chapter, but the importance is not to be ignored. By covering this topic, now we can feel like we can start writing a server application that works safely!

## Exercise
If you see the code for `/error/client/{code}` endpoint, you may see that it handles only the code of "400". We want to change it as follows:

1. Don't you think that the path `code` being used as a string a bit weird? And what if a user hits an path like `/error/client/foo`? Does it make sense to have a value like `"foo"` to be used as a server status code? Let's parse `code` as an integer value, and use it instead of a literal string value of `code`. Of course, you must handle the parsing error!

2. Add a few more client-side error handling lines to the `switch` clause above! You can reference [HTTP Cats](https://http.cat/), a funny and cute website collecting most(if not all) of the HTTP error codes with funny cat images attached. 
