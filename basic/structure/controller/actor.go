package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"structure/service"
)

type ActorController struct {
	path     string
	handlers []ControllerHandler
}

func (c *ActorController) Path() string {
	return c.path
}

func (c *ActorController) Handlers() []ControllerHandler {
	return c.handlers
}

// NewActorController generates a new NewActorController instance.
//
// Inside the `handler` field, we specify the individual API handlers.
func NewActorController() *ActorController {
	return &ActorController{
		path: "/actors",
		handlers: []ControllerHandler{
			{
				Path:   "/",
				Method: http.MethodGet,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "GET /actors/")
				},
			},
			{
				Path:   "/{id}",
				Method: http.MethodGet,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					actorId := r.PathValue("id")
					fmt.Fprintln(w, "GET /actors/{id}")
				},
			},
			{
				Path:   "/",
				Method: http.MethodPost,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					defer r.Body.Close()

					dto := service.ActorUpsertDto{}
					if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
						http.Error(w, "json body parsing error", http.StatusBadRequest)
					}
					fmt.Fprintln(w, "POST /actors/")
				},
			},
			{
				Path:   "/{id}",
				Method: http.MethodPatch,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					// data parsing
					actorId := r.PathValue("id")
					defer r.Body.Close()

					dto := service.ActorUpsertDto{}
					if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
						http.Error(w, "json body parsing error", http.StatusBadRequest)
					}

					// here goes the service layer
					// here we hand over the dto to the service layer
					fmt.Fprintln(w, "PATCH /actors/{id}")
				},
			},
			{
				Path:   "/{id}",
				Method: http.MethodDelete,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					actorId := r.PathValue("id")
					fmt.Fprintln(w, "Delete /actors/{id}")
				},
			},
		},
	}
}
