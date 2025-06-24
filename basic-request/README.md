# Basic Request
In the previous chapter we covered how to write HTTP responses in the `net/http` package. Now we can *write and see the results* right away, it's time to *read* requests and implement our own logic to handle the requests. 

**REMARK** There is a reason why I have chosen to write the response chapter **before** the request chapter. Even if we know how to read things from requests, it is of no use if we cannot check the result that our logic produces from the information that the requests give. As this book is written for programmer's fun, it is **very** important to check the results immediately and feel rewarded.

As usual, run `go run .` on your terminal in this directory and check responses against endpoints:
- `/method`: echoes the method of the request
- `/params/{path}?{query_key1}={query_value1}&{query_key2}={query_value2}`: echoes the parameters of the request; attach extra paths like `/params/wolfpack?phil=bradley&alan=zach`

## A brief review on HTTP requests
Let's say, that writing a HTTP response using `net/http` is relatively easy, as it is the data that all we have to take care of. However, when it comes to reading a HTTP request, there is a ton of information and we need to know important pieces of it. So before we dive in reading HTTP requests, let's briefly skim through the basic parts of HTTP requests that we need to read from them.

For more detailed information, I recommend consulting the [MDN guide on HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Overview).

### Method
The method in a HTTP request is the purpose of it. If it is `GET`, the request want to *read* some information from the server. If it is `POST`, for example, the request asks the server to *create* a new piece of information on the server. 

There are several HTTP methods but throughout the book we will talk about the following methods mainly:

- `GET`: read information(Read)
- `POST`: create information(Create)
- `Patch`: update information(Update)
- `Delete`: delete information(Delte)

You may have heard the acronym "CRUD" a few times. It stands for these four basic operations that are necessary for a user-interacting application to perform. 

### Path
The path of a request is the target, destination of it. It looks like a POSIX path that is used in Linux and macOS. For example, if the path is `/movie/hangover` and the method is `GET`, the request wants to read some information assigned to the path `/movie/hangover`. Of course, the actual information is arbitrary and depends on how the server provides the data in turn. 

### Header
Headers are metadata of HTTP requests. HTTP responses also have headers, but we haven't covered in the previous chapter, as writing custom headers is a bit advanced topic. On the other hand, it is necessary to read information from request headers, because some headers tell significant information about the client, and the server needs to treat it with much care. 

There are lots of headers. For example, `Content-Type` header tells the data type that the client attaches to the request(usually requests of `POST` or `PATCH`). 

### Body
The body part of a HTTP request is where the data that the client wants to send to the server. If a blog user wants to create a blog article, then the browser sends a `POST` request to the blog server with its body containing text data of the article. 

As you notice, `GET` and `Delete` don't need body and in most cases, the server would ignore the content from the body of those requests, even if it is not empty.

## Reading from `http.Request`
It has been a real brief summary of the essential parts of HTTP requests. Now, let's write some Go code and read information from external requests. 

To check results for yourself, use tools like `curl` to check out the responses from the endpoints.

### Method
In `net/http`, reading the method of a request is very straightforward:

```go
mux.HandleFunc("/method", func(w http.ResponseWriter, r *http.Request) {
    var methodPurpose string

    switch r.Method {
    case http.MethodGet:
        methodPurpose = "READ"
    case http.MethodPost:
        methodPurpose = "CREATE"
    case http.MethodPatch:
        methodPurpose = "PATCH"
    case http.MethodDelete:
        methodPurpose = "DELETE"
    default:
        methodPurpose = "Not handled in this echo endpoint"
    }

    w.Write([]byte(r.Method + "; " + methodPurpose + "\n"))
})
```

Use `r.Method` to check the type of the method. In `net/http`, the names of the methods are defined as constants, so you don't have to convert the entire string of `r.Method` into either lowercase or uppercase and compare it with your own custom string, but simply compare with `http.Method{Name}`. The code example above shows how to check the type of the method efficiently, if you are using `net/http`.  

### Path and Query parameters

### Header

### Body

### Form

## Conclusion

## Exercise
