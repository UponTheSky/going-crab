package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/cors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		switch r.Method {
		// preflight
		case http.MethodOptions:
			w.Header().Add("Access-Control-Allow-Methods", http.MethodGet)
			w.Header().Add("Access-Control-Allow-Methods", http.MethodOptions)
			w.Header().Add("Access-Control-Allow-Methods", http.MethodPatch)
			w.Header().Add("Access-Control-Allow-Headers", "X-My-Header")
			w.WriteHeader(http.StatusNoContent)
		default:
			fmt.Fprintln(w, "testing cors successfully!")
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
