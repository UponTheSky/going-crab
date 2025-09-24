package service

import (
	"structure/service/dto"
	"structure/service/schema"
)

type actor = schema.Actor

type IActorService interface {
	ReadAll() ([]actor, error)

	Read(id string) (actor, error)

	Create(dto dto.ActorUpsertDto) (actor, error)

	Update(id string, dto dto.ActorUpsertDto) (actor, error)

	Delete(id string) (actor, error)
}
