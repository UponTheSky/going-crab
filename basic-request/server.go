package main

import (
	"log"
	"net"
	"net/http"
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
}
