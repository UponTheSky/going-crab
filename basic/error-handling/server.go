package main

import (
	"errors"
	"log"
	"net"
	"net/http"
)

func main() {
	// listener
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Cannot open a listener socket: %v", err)
	}

	defer listener.Close()

	// mux
	mux := http.NewServeMux()
	registerEndpoints(mux)

}

func registerEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/error/client/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		switch code {
		case "400":
			http.Error(w, "Bad request", http.StatusBadRequest)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

	mux.HandleFunc("/error/server", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("some kind of error you met while executing the internal logic")

		http.Error(w, err.Error(), http.StatusInternalServerError)
	})
}
