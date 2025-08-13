package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/cookie", func(w http.ResponseWriter, r *http.Request) {
		// set cookie here
		kenCookie := &http.Cookie{
			Name:  "hangover",
			Value: "Ken Jeong",
		}

		zachCookie := &http.Cookie{
			Name:  "hangover",
			Value: "Zach Galifianakis",
		}

		http.SetCookie(w, kenCookie)
		http.SetCookie(w, zachCookie)

		fmt.Fprintln(w, "cookies are set!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
