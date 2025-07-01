package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	// mux
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Listen and Serve!")
	})

	// listener
	listener, err := net.Listen("tcp", ":8080") // if you need to pass context, use ListeConfig.Listen

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	// server
	server := &http.Server{Handler: mux} // could add more configurations - in the later chapters
	defer server.Close()

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
