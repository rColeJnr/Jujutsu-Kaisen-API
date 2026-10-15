package models

type Character struct {
	ID               int64             `json:"id"`
	Name             string            `json:"name"`
	About            string            `json:"about"`
	Alias            []string          `json:"alias"`
	Species          string            `json:"species"`
	Birthday         string            `json:"birthday"`
	Age              int               `json:"age"`
	Gender           string            `json:"gender"`
	Height           int               `json:"height_cm"`
	Weight           float64           `json:"weight_kg"`
	Clan             string            `json:"clan"`
	CursedTechniques []CursedTechnique `json:"cursed_techniques,omitempty"`
	Domain           *DomainExpansion  `json:"domain,omitempty"`
	HairColor        string            `json:"hair_color"`
	EyeColor         string            `json:"eye_color"`
	Status           string            `json:"status"`
	Relatives        []Relative        `json:"relatives,omitempty"`
	Occupation       []Occupation      `json:"occupation,omitempty"`
	Affiliation      []Affiliation     `json:"affiliation,omitempty"`
	Grade            string            `json:"grade"`
	FirstAppearance  *Episode          `json:"first_appearance,omitempty"`
	Images           []Image           `json:"images,omitempty"`
}

type Clan struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Rank        int    `json:"rank"`
}

type DomainExpansion struct {
	ID          int64  `json:"id"`
	CharacterID int64  `json:"character_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CursedTechnique struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Image struct {
	URL string `json:"url"`
}

type Affiliation struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Occupation struct {
	Location string `json:"location"`
	Title    string `json:"title"`
}

type Relative struct {
	CharacterID int64  `json:"character_id"`
	Name        string `json:"name"`
	Relation    string `json:"relation"`
}
