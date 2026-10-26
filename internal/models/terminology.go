package models

type Terminology struct {
	ID      int64  `json:"id"`
	Term    string `json:"term"`
	Meaning string `json:"meaning"`
}

type CreateTerminologyRequest struct {
	Term    string `json:"term"`
	Meaning string `json:"meaning"`
}
type UpdateTerminologyRequest struct {
	Term    *string `json:"term"`
	Meaning *string `json:"meaning"`
}

type TerminologyFilter struct {
	Term string
	Pagination
}
