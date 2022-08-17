package collection_repository

import (
	"time"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

type Collection struct {
	Id        uuid.UUID  `gorm:"primary_key;column:id"`
	Name      string     `gorm:"column:name"`
	AuthorId  uuid.UUID  `gorm:"column:author_id"`
	Topics    []string   `gorm:"column:topics"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (c *Collection) ToEntity() *entity.Collection {
	return &entity.Collection{
		Id:     c.Id,
		Name:   c.Name,
		Topics: c.Topics,
	}
}

type CollectionUserProgress struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CollectionId uuid.UUID `json:"collection_id,omitempty"`
	UserId       uuid.UUID `json:"user_id,omitempty"`
	Mastered     uint32    `json:"mastered,omitempty"`
	Reviewing    uint32    `json:"reviewing,omitempty"`
	Learning     uint32    `json:"learning,omitempty"`
}

func (c *CollectionUserProgress) ToEntity() *entity.CollectionUserProgress {
	return &entity.CollectionUserProgress{
		Id:           c.Id,
		CollectionId: c.CollectionId,
		UserId:       c.UserId,
		Mastered:     c.Mastered,
		Reviewing:    c.Reviewing,
		Learning:     c.Learning,
	}
}

type CollectionUserMetrics struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CollectionId uuid.UUID `json:"collection_id,omitempty"`
	Liked        bool      `json:"liked,omitempty"`
	Disliked     bool      `json:"disliked,omitempty"`
	Viewed       bool      `json:"viewed,omitempty"`
}

func (c *CollectionUserMetrics) ToEntity() *entity.CollectionUserMetrics {
	return &entity.CollectionUserMetrics{
		Id:           c.Id,
		CollectionId: c.CollectionId,
		Liked:        c.Liked,
		Disliked:     c.Disliked,
		Viewed:       c.Viewed,
	}
}

type CollectionMetrics struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CollectionId uuid.UUID `json:"collection_id,omitempty"`
	Likes        uint32    `json:"likes,omitempty"`
	Dislikes     uint32    `json:"dislikes,omitempty"`
	Views        uint32    `json:"views,omitempty"`
}

func (c *CollectionMetrics) ToEntity() *entity.CollectionMetrics {
	return &entity.CollectionMetrics{
		Id:           c.Id,
		CollectionId: c.CollectionId,
		Likes:        c.Likes,
		Dislikes:     c.Dislikes,
		Views:        c.Views,
	}
}

type CollectionCards struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CardId       uuid.UUID `json:"card_id,omitempty"`
	CollectionId uuid.UUID `json:"collection_id,omitempty"`
}

func (c *CollectionCards) ToEntity() *entity.CollectionCards {
	return &entity.CollectionCards{
		Id:           c.Id,
		CardId:       c.CardId,
		CollectionId: c.CollectionId,
	}
}
