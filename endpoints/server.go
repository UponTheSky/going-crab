package main

import (
	"fmt"
	"log"
	"net/http"
)

type LanguageHandler struct {
	language string
}

// This type becomes the Handler type
func (l *LanguageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %v!", l.language)
}

func main() {
	mux := http.NewServeMux()

	// using HandleFunc
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Server!")
	})

	// using Handle, by separately defining a struct implementing the Handler type
	golang := &LanguageHandler{language: "go"}
	rust := &LanguageHandler{language: "rust"}
	swift := &LanguageHandler{language: "swift"}

	mux.Handle("/hello/go/{$}", golang)
	mux.Handle("/hello/rust/{$}", rust)
	mux.Handle("/hello/swift/{$}", swift)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

	// question: what happens when you replace "mux" with "golang", a Handler
	// if err := http.ListenAndServe(":8080", mux); err != nil {
	// 	log.Fatal(err)
	// }
}
