package controller

import "net/http"

type Controller interface {
	Path() string                  // prefix path
	Handlers() []ControllerHandler // individual API handlers
}

type ControllerHandler struct {
	Path        string           // individual subpath
	Method      string           // allowed method for the subpath
	HandlerFunc http.HandlerFunc // individual http handler
}
