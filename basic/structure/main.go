package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"structure/controller"
	"structure/repository"
	"structure/service"
)

func registerController(mux *http.ServeMux, c controller.Controller) {
	for _, handler := range c.Handlers() {
		pattern := fmt.Sprintf("%v %v", handler.Method, filepath.Join(c.Path(), handler.Path))

		mux.Handle(pattern, handler.HandlerFunc)
	}
}

func main() {
	// mux
	mux := http.NewServeMux()

	repository := repository.NewActorRepository(repository.MockDB)
	actorService := service.NewActorService(repository)
	actorController := controller.NewActorController(actorService)
	registerController(mux, actorController)

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
