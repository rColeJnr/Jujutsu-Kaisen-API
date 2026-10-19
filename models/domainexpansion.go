package models

type DomainExpansion struct {
	ID          int64  `json:"id"`
	CharacterID int64  `json:"character_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateDomainExpansionRequest struct {
	CharacterID int64
	Name        string
	Description string
}

type UpdateDomainExpansionRequest struct {
	Name        *string
	Description *string
}
