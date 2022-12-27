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
	CollectionId uuid.UUID `json:"collectionId,omitempty"`
	UserId       uuid.UUID `json:"userId,omitempty"`
	Mastered     uint32    `json:"mastered,omitempty"`
	Reviewing    uint32    `json:"reviewing,omitempty"`
	Learning     uint32    `json:"learning,omitempty"`
}
