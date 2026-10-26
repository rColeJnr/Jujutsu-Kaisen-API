package models

type CursedTechnique struct {
	ID          int64  `json:"id"`
	CharacterID int64  `json:"character_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type CreateCursedTechniqueRequest struct {
	CharacterID int64  `json:"character_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type UpdateCursedTechniqueRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
}

type CursedTechniqueFilter struct {
	Name string
	Pagination
}
