package repository

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"structure/service/dto"
)

var MockDB = map[string]actor{
	"chow": {
		Id:   "chow",
		Name: "Ken Jeong",
		Role: "Lesley Chow",
	},
	"alan": {
		Id:   "alan",
		Name: "Zach Galifianakis",
		Role: "Alan Garner",
	},
}

type SimpleDB = map[string]actor

type ActorRepositoryImpl struct {
	db SimpleDB // usually we have *sql.Conn here
}

func (r *ActorRepositoryImpl) ReadAll() ([]actor, error) {
	return slices.Collect(maps.Values(r.db)), nil
}

func (s *ActorRepositoryImpl) Read(id string) (actor, error) {
	if actor, ok := s.db[id]; ok {
		return actor, nil
	}

	return actor{}, fmt.Errorf("actor with id %v not found", id)
}

func (s *ActorRepositoryImpl) Create(dto dto.ActorUpsertDto) (actor, error) {
	// validation
	if dto.Name == "" || dto.Role == "" {
		return actor{}, errors.New("the name and role must be provided")
	}

	// This could be random, like uuid
	id := strings.Split(dto.Name, " ")[0]
	newActor := actor{
		Id:   id,
		Name: dto.Name,
		Role: dto.Role,
	}

	s.db[id] = newActor

	return newActor, nil
}

func (s *ActorRepositoryImpl) Update(id string, dto dto.ActorUpsertDto) (actor, error) {
	// validation
	if _, ok := s.db[id]; !ok {
		return actor{}, fmt.Errorf("actor with id %v not found", id)
	}

	if dto.Name == "" || dto.Role == "" {
		return actor{}, errors.New("the name and role must be provided")
	}

	updatedActor := actor{Id: id, Name: dto.Name, Role: dto.Role}
	s.db[id] = updatedActor

	return updatedActor, nil
}

func (s *ActorRepositoryImpl) Delete(id string) (actor, error) {
	// validation
	if _, ok := s.db[id]; !ok {
		return actor{}, fmt.Errorf("actor with id %v not found", id)
	}

	deletedActor := s.db[id]
	delete(s.db, id)

	return deletedActor, nil
}

func NewActorRepository(db SimpleDB) *ActorRepositoryImpl {
	return &ActorRepositoryImpl{db: db}
}
