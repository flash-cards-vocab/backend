package entity

import (
	"github.com/google/uuid"
)

type CollectionFullUserMetricsResponse struct {
	CollectionId uuid.UUID `json:"collectionId,omitempty"`
	Likes        uint32    `json:"likes"`
	Dislikes     uint32    `json:"dislikes"`
	Views        uint32    `json:"views"`
	UserId       uuid.UUID `json:"userId"`
	Liked        bool      `json:"liked"`
	Disliked     bool      `json:"disliked"`
	Viewed       bool      `json:"viewed"`
	Starred      bool      `json:"starred"`
}

type CollectionMetrics struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CollectionId uuid.UUID `json:"collectionId"`
	Likes        uint32    `json:"likes"`
	Dislikes     uint32    `json:"dislikes"`
	Views        uint32    `json:"views"`
}

type CollectionUserMetrics struct {
	Id           uuid.UUID `json:"id,omitempty"`
	UserId       uuid.UUID `json:"userId"`
	CollectionId uuid.UUID `json:"collectionId"`
	Liked        bool      `json:"liked"`
	Disliked     bool      `json:"disliked"`
	Viewed       bool      `json:"viewed"`
	Starred      bool      `json:"starred"`
}
