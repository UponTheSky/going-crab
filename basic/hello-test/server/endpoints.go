package server

import (
	"hello-test/app"
	"io"
	"log"
	"net/http"
)

func CapitalizeHandler(logger *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType, ok := r.Header["Content-Type"]

		if !ok || len(contentType) == 0 || (contentType[0] != "text/plain") {
			msg := "the header Content-Type is not text/plain"
			logger.Println(msg)

			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		text, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			logger.Printf("error while reading r.Body: %v", err)

			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		body := app.Capitalize(string(text))
		w.Write([]byte(body))
	}
}
