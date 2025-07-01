package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	// register endpoints
	mux := http.NewServeMux()
	registerEndpoints(mux)

	// create server
	server := &http.Server{Handler: mux}
	defer server.Close()

	// create listener
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	// run the server
	err = server.Serve(listener)

	log.Fatal(err)
}

func registerEndpoints(mux *http.ServeMux) {
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

	mux.HandleFunc("/params/{pathParam}", func(w http.ResponseWriter, r *http.Request) {
		// path params
		pathParam := r.PathValue("pathParam")

		// queries
		queryParams := r.URL.Query()
		queryStringBuilder := &strings.Builder{}

		for key, value := range queryParams {
			queryStringBuilder.WriteString("(key=" + key + ",value=" + value[0] + ") ")
		}

		w.Write([]byte("path: " + pathParam + "\nquery: " + queryStringBuilder.String() + "\n"))
	})

	mux.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header

		headerStringBuilder := &strings.Builder{}

		for key, value := range headers {
			headerStringBuilder.WriteString("(key=" + key + ", value=" + value[0] + ") ")
		}

		w.Write([]byte("headers: " + headerStringBuilder.String() + "\n"))
	})

	mux.HandleFunc("/body", func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 256)
		n, err := r.Body.Read(buf)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		w.Write(buf[0:n])
		w.Write([]byte("\n"))

	})

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
}
