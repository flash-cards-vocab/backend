package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserCollectionResponse struct {
	Id               uuid.UUID `json:"id,omitempty"`
	Name             string    `json:"name"`
	AuthorName       string    `json:"authorName"`
	Topics           []string  `json:"topics"`
	Starred          bool      `json:"starred"`
	Likes            uint32    `json:"likes"`
	Dislikes         uint32    `json:"dislikes"`
	Views            uint32    `json:"views"`
	Mastered         uint32    `json:"mastered"`
	Reviewing        uint32    `json:"reviewing"`
	Learning         uint32    `json:"learning"`
	IsLikedByUser    bool      `json:"isLikedByUser"`
	IsDislikedByUser bool      `json:"isDislikedByUser"`
	IsViewedByUser   bool      `json:"isViewedByUser"`
	TotalCards       int       `json:"totalCards"`
	CreatedDate      string    `json:"createdDate"`
}

type CollectionPreviewResponse struct {
	Id               uuid.UUID `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Topics           []string  `json:"topics"`
	Starred          bool      `json:"starred"`
	Likes            uint32    `json:"likes"`
	Dislikes         uint32    `json:"dislikes"`
	Views            uint32    `json:"views"`
	IsLikedByUser    bool      `json:"isLikedByUser"`
	IsDislikedByUser bool      `json:"isDislikedByUser"`
	IsViewedByUser   bool      `json:"isViewedByUser"`
	TotalCards       int       `json:"totalCards"`
}

type Collection struct {
	Id        uuid.UUID  `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Topics    []string   `json:"topics,omitempty"`
	AuthorId  uuid.UUID  `json:"author_id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateCollectionRequest struct {
	Name   string   `json:"name,omitempty"`
	Topics []string `json:"topics,omitempty"`
	Cards  []*Card  `json:"cards,omitempty"`
}

type GetCollectionWithCardsResponse struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Mastered   uint32    `json:"mastered"`
	Reviewing  uint32    `json:"reviewing"`
	Learning   uint32    `json:"learning"`
	TotalCards int       `json:"totalCards,omitempty"`
	Cards      []*Card   `json:"cards,omitempty"`
}
