package dto

// ActorUpsertDto contains necessary information for creating and updating an Actor object
type ActorUpsertDto struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
