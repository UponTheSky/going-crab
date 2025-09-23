package main

import (
	"structure/controller"
)

type App struct {
	controllers []controller.Controller
}

func (a *App) Run(servera any) {

}

// TODO: add middlewares
// for individual controllers, we need to apply the visitor pattern to add the middlewares
// for individual API endpoints, we simply add middlewares directly
func NewApp(controllers []controller.Controller) *App {
	return &App{controllers: controllers}
}
