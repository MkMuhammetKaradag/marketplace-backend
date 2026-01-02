package domain

import "github.com/google/uuid"

type Category struct {
	ID          uuid.UUID `json:"id"`
	ParentID    uuid.UUID `json:"parent_id,omitempty"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}
