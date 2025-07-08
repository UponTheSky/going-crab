# Test Your Server Automatically
By far, we have checked out our servers through manual testing, using tools such as `curl` or [Postman](https://www.postman.com/). Manual testing is good in its own way, because it is simple to make a test and obtain the result immediately. Thus it has suited our needs over the course of the past chapters. 

However, as we will build servers with more complex structures, manual testing is very limited in that:
- it is only possible to check out the entire behavior of the server rather than a piece of the internal business logic, and
- it takes too much time and effort to replicate tests whenever there are changes in the implementation.

Thus, we would like to find someone to do these hassle tests on behalf; for sure AI(or LLM) could be an answer and there are already movements towards automating tests with AI. However, in this chapter we will test code on our own. I believe by writing test code, you can understand the requirements that you want to convert to computer logic deeper. 

Here we will only cover basic part of writing test code in Go. For learning a more hands-on approach, I highly recommend the book ["Learn Go with Tests"](https://quii.gitbook.io/learn-go-with-tests) by [Chris James](https://github.com/quii).

## Project Structure
Let's say that we have an API endpoint `/capitalize` on our server such that any `POST` request to it with a text body returns the text with all of its letters capitalized. So 

```sh
curl -X POST localhost:8080/capitalize -d "Ken Jeong"
```

would return:

```sh
"KEN JEONG"
```

As we have done in the previous chapters, we will do a very minimalistic setup for our server. This time, for the purpose of making test, we will separate the business logic from the package `main` which is for execution. 

The structure of the example project for this chapter would be as follows:

```sh
├── app
│   ├── capitalize_test.go
│   └── capitalize.go
├── go.mod
├── main.go
├── README.md
└── server
    ├── endpoints_test.go
    ├── endpoints.go
    └── server.go
```

- `/app`: the directory for the files about the business logic. The suffix `_test` for `{file}_test.go` means it has the test logic for the file `{file}.go`.
- `/server`: the directory for setting up the server. `endpoints_test` will test the "logic" in the `endpoints.go` file. Wait, is there something to be called *logic* for just running up the server with bunch of endpoints? We'll see...

The execution will occur in `main.go` but it would not have any significant information. So our forus in this chapter will be `/app` and `/server`.

## Writing Tests
The term *unit test* usually refers to testing a separate logical module or a function without dependencies if possible. You can imagine a complex machine consisting of small gears and cogs. Unit test is about testing a small gear; even one single small gear could cause a trouble for the entire machine, so we have to make sure every single gear works okay. 

On the other hand, *integration test* is about testing how two or more gears work together. Once each individual gear is fine, then we check out how multiple gears rotate together. 

There are several kinds of tests, but in this chapter we will mainly cover these two. 

### Unit Tests
In our small project, the business logic is in `app/capitalize.go`:

```go
import "strings"

func Capitalize(text string) string {
	return strings.ToUpper(text)
}
```
(I know, it is a one-liner simple function, but for the sake of easy explanation)

To test this function, we write a test in `app/capitalize_test.go`:

```go
package app

import "testing"

func TestEmpty(t *testing.T) {
    // asset
	input := ""
	want := ""

    // act
	got := Capitalize(input)

    // assert
	if got != want {
		t.Errorf("Test Empty - got=%v, want=%v", got, want)
	}
}

func TestText(t *testing.T) {
    // asset
	input := "Ken Jeong!"
	want := "KEN JEONG!"

    // act
	got := Capitalize(input)

    // assert 
	if got != want {
		t.Errorf("Test Empty - got=%v, want=%v", got, want)
	}
}
```

It is a very simple logic:
1. You start with preparing the input and the result you want(*asset*, *given*).
2. Compute the result from the asset prepared in the previous step(*action*, *when*).
3. Compare the result and the expected value(*assert*, *then*).

Some people call these three steps as *Asset-Act-Assert*, and others call as *Given-When-Then*. Whatever the calling terms they are, the processes are pretty straightforward. 

Now, run the test by executing the following command `go test ./app -v` and see the results like this:

```sh
=== RUN   TestEmpty
--- PASS: TestEmpty (0.00s)
=== RUN   TestText
--- PASS: TestText (0.00s)
PASS
ok      hello-test/app  0.165s
```

### Integration Tests

#### Server Setup
After testing individual modules, we need to see if those modules work coherently in the same system. In our project, we want to make sure the endpoint `/capitalize` works correctly.

First of all, just like we have done over the past chapters, set up a TCP server in a simple way at the port `8080`, in `server/server.go`:

```go
package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type Server struct {
	protocol string
	port     int
	mux      *http.ServeMux
	logger   *log.Logger
}

func New(protocol string, port int) *Server {
	mux := http.NewServeMux()
	logger := log.New(
		os.Stderr,
		"[Going Crab] ",
		log.LstdFlags|log.Lshortfile|log.LUTC,
	)

	return &Server{
		protocol: protocol,
		port:     port,
		mux:      mux,
		logger:   logger,
	}
}

func (ws *Server) RegisterEndpoint(path string, handler http.Handler) {
	ws.mux.Handle(path, handler)
}

func (ws *Server) Run() error {
	srv := &http.Server{}
	srv.Addr = fmt.Sprintf(":%v", ws.port)
	srv.Handler = ws.mux

	listener, err := net.Listen(ws.protocol, srv.Addr)

	if err != nil {
		return err // unexpected error
	}

	defer listener.Close()

	return srv.Serve(listener) // server.Serve always returns an error
}
```
Basically it is a collection of necessary procedures for running a TCP server, hence there is not much logic to be tested. If you *dare* to test those functions, consult the source code of the package where the functions come from. The Go team already has tested them thoroughly. By the way, the code would be straightforward if you have followed the book in order.  

On the other hand, we separately implement the logic for handling HTTP requests & responses in `endpoints.go`. Note that we inject the logger here:

```go
package server

import (
	"hello-test/app"
	"io"
	"log"
	"net/http"
)

func CapitalizeHandler(logger *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType, ok := r.Header["Content-Type"]

		if !ok || len(contentType) == 0 || (contentType[0] != "text/plain") {
			msg := "the header Content-Type is not text/plain"
			logger.Println(msg)

			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		text, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			logger.Printf("error while reading r.Body: %v", err)

			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		body := app.Capitalize(string(text))
		w.Write([]byte(body))
	}
}
```

Thus, our main focus for running an integration test for the endpoint `/capitalize` will be on `endpoints_test.go`. Now I would like to introduce a new package [`http/net/httptest`](https://pkg.go.dev/net/http/httptest) that will help us writing HTTP-server related test code in a fairly simple way. The package itself is a kind of small giftbox for those who love to use `net/http` package.

#### Writing Test with TestServer
First, we will test a server that is actually running. You may think we need to run the actual instance of the server we are testing, and make HTTP client requests just like we do when making external API calls. However, `http/httptest` package provides a set of convenient tools so we don't have to.

Our strategy here also follows the same *Asset-Act-Assert* pattern:
1. Asset: Start a server using `httptest.NewServer`, without changing the default configs(we may need to do it in the future chapters).
2. Act: Make a client call using `*Server.Client()` with input text value attached as its body.
3. Assert: compare the result from the response and the expected result. 

Starting a new server using `httptest` is fairly simple, but we need to prepare our mutliplexer where all the endpoints are registered:

```go
package server

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEndpointWithServer(t *testing.T) {
	// asset
	mux := http.NewServeMux()
	logger := log.New(os.Stderr, "", log.LstdFlags)

	mux.HandleFunc("/capitalize", CapitalizeHandler(logger))

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	input := "Bradley Cooper"
	want := "BRADLEY COOPER"

	// [...]
}
```
It seems a little bit daunting compared to the unit tests covered above, but still nothing complex here. We generate a mux and logger, and then provide them to [`httptest.NewServer`](https://pkg.go.dev/net/http/httptest#NewServer) to create an run a [test server](https://pkg.go.dev/net/http/httptest#Server).

Next, we make a POST call using the generated test `Server`'s [`Client`](https://pkg.go.dev/net/http/httptest#Server.Client):

```go
func TestEndpointWithServer(t *testing.T) {
	// asset
	// [...]

	// act
	client := testServer.Client()
	body := strings.NewReader(input)

	response, err := client.Post(testServer.URL+"/capitalize", "text/plain", body)

	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	// [...]
}
```
Since one of our requirements for the endpoint `/capitalize` is to have `Content-Type` header's value of `text/plain`, we specify it when making a POST call. 

Finally, check if we have got the things we expected:

```go
func TestEndpointWithServer(t *testing.T) {
	// [...]

	// assert
	if response.StatusCode != http.StatusOK {
		t.Errorf("the server response status code: got=%v, want=%v", response.StatusCode, http.StatusOK)
	}

	respBody := response.Body
	defer respBody.Close()

	gotByte, err := io.ReadAll(respBody)

	if err != nil {
		t.Fatalf("unexpected response read error: %v", err)
	}

	got := string(gotByte)

	if got != want {
		t.Errorf("Test /capitalize endpoint: got=%v, want=%v", got, want)
	}
}
```
First, we check the status code of the response, so that whether the POST request itself has had no problem. Then we check the data we received from the server, and compare it with the value we expected.

Try `go test ./server -v` and see what you get on your terminal.

#### Writing Test without TestServer
However, we don't really need to run the test server itself; if what we test is only the handlers, then why don't we inject mock requests and responses to the handlers? `http/httptest` package provides useful tools in this case as well.

For the asset preparation step, we prepare the following:
- handler(s)
- data we want to compare(input, expected value)
- a `*http.Request` generated from [`httptest.NewRequest`](https://pkg.go.dev/net/http/httptest#NewRequest)
- a recorder mocking our response, generated from [`httptest.NewRecorder`](https://pkg.go.dev/net/http/httptest#NewRecorder)

```go
func TestEndpointNoServer(t *testing.T) {
	// asset

	// handler
	logger := log.New(os.Stderr, "", log.LstdFlags)
	handler := CapitalizeHandler(logger)

	// data
	input := "Bradley Cooper"
	want := "BRADLEY COOPER"

	// *http.Request
	body := strings.NewReader(input)
	request := httptest.NewRequest(http.MethodPost, "/capitalize", body)
	request.Header.Set("Content-Type", "text/plain")

	// *httptest.ResponseRecorder
	rw := httptest.NewRecorder()

	// [...]
}
```

The next step is very simple. You just need to invoke the handler

```go
func TestEndpointNoServer(t *testing.T) {
	// [...]

	// act
	handler(rw, request)

	// [...]
}
```

Finally, we check whether the response was okay, and the received result is as expected:

```go
func TestEndpointNoServer(t *testing.T) {
	// [...]

	// assert
	response := rw.Result()

	if response.StatusCode != http.StatusOK {
		t.Errorf("status code error - want=%v, got=%v", http.StatusOK, response.StatusCode)
		return
	}

	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	got := string(resBody)

	if got != want {
		t.Errorf("different return value - want=%v, got=%v", want, got)
		return
	}
}
```

The code is much simpler than running a test server. Unless you have an endpoint closely related to how the server is configured, we recommend using this second hack. 

## Conclusion
Since we can test through code, we are set free from those hassle days of manual testing. However, code is a liability. If we change the behavior of a certain part of our logic, we also have to change its corresponding test code. Sometimes, badly written test code is worse than not writing any test code. You should be aware of maintaining test code in clear and concise ways.

## Exercise
In `/server/endpoints_test.go`, we only covered *good paths*. That is, we only test the case when the client request is normal. However, the world is full of evils - *bad paths*. In this exercise, you need to add a test function for one of them.

Add another test function to `endpoints_test.go`, where a given request has the header `Content-Type` different from `text/plain`. You may use various techniques here - you may directly send a JSON body, make a request with a form, or directly set the value of the header. Oh, and please use the technique of not running a test server. 
