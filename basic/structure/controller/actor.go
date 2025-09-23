package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"structure/service"
	"structure/service/dto"
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
func NewActorController(actorService service.IActorService) *ActorController {
	return &ActorController{
		path: "/actors",
		handlers: []ControllerHandler{
			{
				Path:   "/",
				Method: http.MethodGet,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					actors, err := actorService.ReadAll()

					if err != nil {
						// log the error
						http.Error(w, "internal server error", http.StatusInternalServerError)
						return
					}

					// write headers first before encoding and returning the data(specific in Go's encode/json)
					w.WriteHeader(http.StatusOK)

					// encode
					if err := json.NewEncoder(w).Encode(&actors); err != nil {
						// json encode error
						http.Error(w, "internal server error", http.StatusInternalServerError)
					}
				},
			},
			{
				Path:   "/{id}",
				Method: http.MethodGet,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					actorId := r.PathValue("id")

					actor, err := actorService.Read(actorId)

					if err != nil {
						// TODO: separately define errors to distinguish between 404 and 401
						// at the moment, assume that it is 404
						http.Error(w, fmt.Sprintf("actor with id %v not found", actorId), http.StatusNotFound)
						return
					}

					w.WriteHeader(http.StatusOK)

					if err := json.NewEncoder(w).Encode(&actor); err != nil {
						// json encode error
						http.Error(w, "internal server error", http.StatusInternalServerError)
					}
				},
			},
			{
				Path:   "/",
				Method: http.MethodPost,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					// parsing the request input
					defer r.Body.Close()

					dto := dto.ActorUpsertDto{}
					if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
						http.Error(w, "json body parsing error", http.StatusBadRequest)
						return
					}

					newActor, err := actorService.Create(dto)

					if err != nil {
						// log the error
						// todo: need to separately define errors for invalid inputs
						// at the moment, we assume that it is 500
						http.Error(w, "internal server error", http.StatusInternalServerError)
						return
					}

					w.WriteHeader(http.StatusCreated)

					if err := json.NewEncoder(w).Encode(&newActor); err != nil {
						// json encode error
						http.Error(w, "internal server error", http.StatusInternalServerError)
					}
				},
			},
			{
				Path:   "/{id}",
				Method: http.MethodPatch,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					// data parsing
					actorId := r.PathValue("id")
					defer r.Body.Close()

					dto := dto.ActorUpsertDto{}
					if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
						http.Error(w, "json body parsing error", http.StatusBadRequest)
						return
					}

					// here goes the service layer
					// here we hand over the dto to the service layer
					updatedActor, err := actorService.Update(actorId, dto)

					if err != nil {
						// log the error
						// todo: need to separately define errors for invalid inputs, or if the actor doesn't exist
						// at the moment, we assume that it is 401
						http.Error(w, "the request contains invalid data", http.StatusBadRequest)
						return
					}

					// write the status header first(specific to when using encode/json's Encoder::Encode())
					w.WriteHeader(http.StatusOK)

					// return the json as response
					if err := json.NewEncoder(w).Encode(&updatedActor); err != nil {
						// json encode error
						http.Error(w, "internal server error", http.StatusInternalServerError)
					}
				},
			},
			{
				Path:   "/{id}",
				Method: http.MethodDelete,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					actorId := r.PathValue("id")

					deletdActor, err := actorService.Delete(actorId)

					if err != nil {
						// log the error
						// TODO: separate the errors
						// for now assume that it is 404
						http.Error(w, fmt.Sprintf("actor of id %v not found", actorId), http.StatusNotFound)
						return
					}

					w.WriteHeader(http.StatusOK)

					if err := json.NewEncoder(w).Encode(&deletdActor); err != nil {
						// json encode error
						http.Error(w, "internal server error", http.StatusInternalServerError)
					}
				},
			},
		},
	}
}
