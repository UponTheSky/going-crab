package service

import (
	"structure/repository"
	"structure/service/dto"
)

type ActorServiceImpl struct {
	actorRepository repository.IActorRepository
}

func (s *ActorServiceImpl) ReadAll() ([]actor, error) {
	return s.actorRepository.ReadAll()
}

func (s *ActorServiceImpl) Read(id string) (actor, error) {
	// TODO: validate the id value
	return s.actorRepository.Read(id)
}

func (s *ActorServiceImpl) Create(dto dto.ActorUpsertDto) (actor, error) {
	// TODO: validate the dto

	// NOTE: starts transaction here if the DB supports transactional operations
	newActor, err := s.actorRepository.Create(dto)

	// NOTE; ends transaction here if the DB supports transactional operations

	return newActor, err
}

func (s *ActorServiceImpl) Update(id string, dto dto.ActorUpsertDto) (actor, error) {
	// TODO: validate the dto

	// NOTE: starts transaction here if the DB supports transactional operations
	updatedActor, err := s.actorRepository.Update(id, dto)

	// NOTE: ends transaction here if the DB supports transactional operations

	return updatedActor, err
}

func (s *ActorServiceImpl) Delete(id string) (actor, error) {
	// TODO: validate the id

	return s.actorRepository.Delete(id)

}

func NewActorService(repository repository.IActorRepository) *ActorServiceImpl {
	return &ActorServiceImpl{actorRepository: repository}
}
