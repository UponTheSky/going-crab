package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"structure/controller"
	"structure/repository"
	"structure/service"
)

type App struct {
	controllers []controller.Controller
}

func (a *App) RegisterController(mux *http.ServeMux) {
	for _, c := range a.controllers {
		for _, handler := range c.Handlers() {
			pattern := fmt.Sprintf("%v %v", handler.Method, filepath.Join(c.Path(), handler.Path))
			mux.Handle(pattern, handler.HandlerFunc)
		}
	}
}

func NewApp() *App {
	repository := repository.NewActorRepository(repository.MockDB)
	actorService := service.NewActorService(repository)

	// TODO: add midlewares for actorController
	actorController := controller.NewActorController(actorService)

	// TODO: add middlewares for the entire app
	return &App{controllers: []controller.Controller{actorController}}
}
