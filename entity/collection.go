package entity

import (
	"github.com/google/uuid"
)

type UserCollectionResponse struct {
	Id               uuid.UUID `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Topics           []string  `json:"topics,omitempty"`
	Likes            uint32    `json:"likes,omitempty"`
	Dislikes         uint32    `json:"dislikes,omitempty"`
	Views            uint32    `json:"views,omitempty"`
	Mastered         uint32    `json:"mastered,omitempty"`
	Reviewing        uint32    `json:"reviewing,omitempty"`
	Learning         uint32    `json:"learning,omitempty"`
	IsLikedByUser    bool      `json:"is_liked_by_user,omitempty"`
	IsDislikedByUser bool      `json:"is_disliked_by_user,omitempty"`
	IsViewedByUser   bool      `json:"is_viewed_by_user,omitempty"`
}

type Collection struct {
	Id     uuid.UUID `json:"id,omitempty"`
	Name   string    `json:"name,omitempty"`
	Topics []string  `json:"topics,omitempty"`
}
