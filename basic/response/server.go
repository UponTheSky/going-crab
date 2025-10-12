package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

type Character struct {
	Name   string `json:"name"`
	Actor  string `json:"actor"`
	Series []int  `json:"series"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/string", stringHandler)
	mux.HandleFunc("/json", jsonHandler)
	mux.HandleFunc("/header", headerHandler)

	staticFileHandler := http.FileServer(http.Dir("./static"))
	mux.Handle("/", staticFileHandler)

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{Handler: mux}

	err = server.Serve(listener)

	log.Fatal(err) // err is always not nil in this case
}

func stringHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a string response\n"))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	chow := Character{
		Name:   "Leslie Chow",
		Actor:  "Ken Jeong",
		Series: []int{1, 2, 3},
	}

	// using json.Marshal
	// chow_json, err := json.Marshal(chow)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// w.Write(chow_json)

	// using Encoder.Encode

	// wrap the writing stream as Encoder
	encoder := json.NewEncoder(w)

	// invoke Encode and directly pass the struct object to the function
	if err := encoder.Encode(chow); err != nil {
		log.Fatal(err)
	}
}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("foo", "foo1") // this will be ignored
	w.Header().Set("foo", "foo2") // this will be set as the header value

	// both will be set
	w.Header().Add("bar", "bar1")
	w.Header().Add("bar", "bar2")

	// set the status code
	w.WriteHeader(http.StatusOK)

	// write the response
	fmt.Fprintln(w, "please check the header")

	w.Header().Set("baz", "this-will-not-showup") // after either w.WriteHeader or w.Write is invoked
}
