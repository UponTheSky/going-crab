package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	logger  *log.Logger
	handler http.Handler
}

func (lm *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log the method of the request, path(endpoint), and the current time
	lm.logger.Println(r.Method, r.URL.Path, time.Now())

	// hands the request over to the given handler
	lm.handler.ServeHTTP(w, r)
}

// func AddLoggingMiddleware(handler http.Handler) http.Handler {
// 	middlware := func(w http.ResponseWriter, r *http.Request) {
// 		// log the method of the request, path(endpoint), and the current time
// 		log.Println(r.Method, r.URL.Path, time.Now())

// 		// hands the request over to the given handler
// 		handler.ServeHTTP(w, r)
// 	}

// 	return http.HandlerFunc(middlware)
// }

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
	// mainHandler := AddLoggingMiddleware(mux)
	mainHandler := &LoggingMiddleware{
		logger:  log.Default(),
		handler: mux,
	}

	// listener
	listener, err := net.Listen("tcp", ":8080") // if you need to pass context, use ListeConfig.Listen

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	// server
	server := &http.Server{Handler: mainHandler}
	defer server.Close()

	err = server.Serve(listener)
	log.Fatal(err)
}
