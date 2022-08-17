package entity

import (
	"github.com/google/uuid"
)

type CollectionMetrics struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CollectionId uuid.UUID `json:"collection_id,omitempty"`
	Likes        uint32    `json:"likes,omitempty"`
	Dislikes     uint32    `json:"dislikes,omitempty"`
	Views        uint32    `json:"views,omitempty"`
}

type CollectionUserMetrics struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CollectionId uuid.UUID `json:"collection_id,omitempty"`
	Liked        bool      `json:"liked,omitempty"`
	Disliked     bool      `json:"disliked,omitempty"`
	Viewed       bool      `json:"viewed,omitempty"`
}
