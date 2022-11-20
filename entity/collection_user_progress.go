package entity

import (
	"github.com/google/uuid"
)

type CollectionUserProgressResponse struct {
	Mastered  uint32 `json:"mastered,omitempty"`
	Reviewing uint32 `json:"reviewing,omitempty"`
	Learning  uint32 `json:"learning,omitempty"`
}

type CollectionUserProgress struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CollectionId uuid.UUID `json:"collection_id,omitempty"`
	UserId       uuid.UUID `json:"user_id,omitempty"`
	Mastered     uint32    `json:"mastered,omitempty"`
	Reviewing    uint32    `json:"reviewing,omitempty"`
	Learning     uint32    `json:"learning,omitempty"`
}
