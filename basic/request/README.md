# Basic Request
In the previous chapter we covered how to write HTTP responses in the `net/http` package. Now we can *write and see the results* right away, it's time to *read* requests and implement our own logic to handle the requests. 

**REMARK** There is a reason why I have chosen to write the response chapter **before** the request chapter. Even if we know how to read things from requests, it is of no use if we cannot check the result that our logic produces from the information that the requests give. As this book is written for programmer's fun, it is **very** important to check the results immediately and feel rewarded.

As usual, run `go run .` on your terminal in this directory and check responses against endpoints:
- `/method`: echoes the method of the request
- `/params/{path}?{query_key1}={query_value1}&{query_key2}={query_value2}`: echoes the parameters of the request; attach extra paths like `/params/wolfpack?phil=bradley&alan=zach`
- `/header`: echoes the list of the headers of the request
- `/body`: echoes the body content of the request
- `/form`: echoes the form content of the request

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

Also note that we have body properties in HTTP responses as well. This is the part in which we wrote texts and JSONs in the previous chapter.

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
Although accessing the entire path of a request is through the [URL object](https://pkg.go.dev/net/url#URL) of the request, accessing the matched path parameters is through `Request.PathValue`:

```go
mux.HandleFunc("/params/{pathParam}", func(w http.ResponseWriter, r *http.Request) {
    // path params
    pathParam := r.PathValue("pathParam")
    // [...]
})
```

Specifying which parameter to match is related with how you design the endpoint string pattern(see the details on the [`http.ServeMux` page](https://pkg.go.dev/net/http#ServeMux)). Here we try to match the parameter right next to `params`.

Accessing query values is through the URL object of the request, `Request.URL.Query`:

```go
mux.HandleFunc("/params/{pathParam}", func(w http.ResponseWriter, r *http.Request) {
    // [...]

    // queries
    queryParams := r.URL.Query()
    queryStringBuilder := &strings.Builder{}

    for key, value := range queryParams {
        queryStringBuilder.WriteString("key: " + key + ", value: " + value[0] + "\n")
    }

    // [...]
})
```

`Request.URL.Query()` returns a [map of string array values](https://pkg.go.dev/net/url#Values), so the above code uses the syntax of a loop over the map to retrieve all the query values.

### Header
Reading the headers of a request is also simple:

```go
mux.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
    headers := r.Header

    headerStringBuilder := &strings.Builder{}

    for key, value := range headers {
        headerStringBuilder.WriteString("(key=" + key + ", value=" + value[0] + ") ")
    }

    w.Write([]byte("headers: " + headerStringBuilder.String() + "\n"))
})
```

Just like the return value of `URL.Query`, `Request.Header` is a simple map with string array values, so the above code tries to catch the very first value of each of the keys. However, note that this doesn't provide all the actual headers that the request has, and some important headers such as `Host` is separated as another property of the request. Check out the [documentation](https://pkg.go.dev/net/http#Request) for details.

### Body
[`Request.Body`](https://pkg.go.dev/net/http#Request) follows [`io.ReadCloser` interface](https://pkg.go.dev/io#ReadCloser). But according to the documentation, we don't have to close it manually. All we have to do is to read data from it, and forget about it.

```go
mux.HandleFunc("/body", func(w http.ResponseWriter, r *http.Request) {
    buf := make([]byte, 256)
    n, err := r.Body.Read(buf) // be careful not to pass "buf" with length zero; `Read` would not read any data. See https://pkg.go.dev/io#Reader.

    if err != nil && err != io.EOF {
        log.Fatal(err)
    }

    w.Write(buf[0:n])
    w.Write([]byte("\n"))
})
```

Here we create our own buffer and pass it to `r.Body.Read`, but we can choose out of many different reading strategies. If the data is not expected to be big, simply reading all at once using [`io.ReadAll`](https://pkg.go.dev/io#ReadAll) is a good strategy. Otherwise, using [`io.Scanner`](https://pkg.go.dev/bufio#Scanner) could reduce syscall call overhead while reading large size of data. 

### Form
[Form](https://datatracker.ietf.org/doc/html/rfc1866#section-8) is a special type of body that is attached to a request of `POST`, `PUT`, or `PATCH`. It is a encoded string data of a collection of (key, value) pairs, which are usually generated from `<form>` HTML tag. 

In `net/http` package, we retrieve such data using `Request.PostForm`:

```go
mux.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        log.Fatal(err)
    }

    form := r.PostForm

    formStringBuilder := &strings.Builder{}

    for key, value := range form {
        formStringBuilder.WriteString("(key=" + key + ", value=" + value[0] + ") ")
    }

    w.Write([]byte("form: " + formStringBuilder.String() + "\n"))
})
```

Similar to headers or query parameters, the retrived values are a map object in Go. Note that you must call [`Request.ParseForm`](https://pkg.go.dev/net/http#Request.ParseForm) in advance, in order to retrive the form values from the given request.

## Conclusion
Although we haven't covered all the important details of reading information from a HTTP request, we can at least write a server implementing simple I/O logics! Let's write a fun server with mock DBs in the next first challenge exercise! But before then, there are three prerequisite chapters that we would be better off covering: error handling, logging, and testing. Of course, we won't get into much of details right at the moment, because we need to enjoy our fruits ASAP! Even if those subjects seem to be a bit boring(and sometimes yes...), you would be greatful for knowing these topics as these reduces your debugging time significantly. Yeah, you've heard, *DEBUGGING*... 

## Exercise
An acute reader may have already noticed that we haven't covered reading requests with JSON bodies. 

1. Use [`bufio.Scanner`](https://pkg.go.dev/bufio#Scanner) API for reading multiple lines of text data from a given request. The text data is expected to be an encoded JSON object.

2. Use [`json.UnMarshal`](https://pkg.go.dev/encoding/json#Unmarshal) to directly parse the entire JSON body. Use [`io.ReadAll`](https://pkg.go.dev/io#ReadAll) for reading the body.

3. Use [`json.Decoder.Decode](https://pkg.go.dev/encoding/json#Decoder.Decode) to handle JSON data in only a few lines! What a life! But can you be sure that 
