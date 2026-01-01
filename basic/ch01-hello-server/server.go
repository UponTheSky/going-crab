package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Server!")
	})

	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		fmt.Println(err)
	}
}
