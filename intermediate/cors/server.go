package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/cors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		fmt.Fprintln(w, "testing cors successfully!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
