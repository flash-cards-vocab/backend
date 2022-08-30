package card_repository

import (
	"time"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

type Card struct {
	Id         uuid.UUID  `gorm:"primary_key;column:id"`
	Word       string     `gorm:"column:word"`
	ImageUrl   string     `gorm:"column:image_url"`
	Definition string     `gorm:"column:definition"`
	Sentence   string     `gorm:"column:sentence"`
	Antonyms   string     `gorm:"column:antonyms"`
	Synonyms   string     `gorm:"column:synonyms"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`
}

func (c *Card) ToEntity() *entity.Card {
	return &entity.Card{
		Id:         c.Id,
		Word:       c.Word,
		ImageUrl:   c.ImageUrl,
		Definition: c.Definition,
		Sentence:   c.Sentence,
		Antonyms:   c.Antonyms,
		Synonyms:   c.Synonyms,
	}
}

type CardUserProgress struct {
	Id            uuid.UUID                   `gorm:"primary_key;column:id"`
	CardId        uuid.UUID                   `gorm:"column:card_id"`
	UserId        uuid.UUID                   `gorm:"column:user_id"`
	Status        entity.CardUserProgressType `gorm:"column:status"`
	LearningCount uint32                      `gorm:"column:learning_count"`
}

func (c *CardUserProgress) ToEntity() *entity.CardUserProgress {
	return &entity.CardUserProgress{
		Id:     c.Id,
		CardId: c.CardId,
		UserId: c.UserId,
		Status: c.Status,
	}
}

type CardMetrics struct {
	Id       uuid.UUID `gorm:"primary_key;column:id"`
	CardId   uuid.UUID `gorm:"column:card_id"`
	Likes    uint32    `gorm:"column:likes"`
	Dislikes uint32    `gorm:"column:dislikes"`
}

func (c *CardMetrics) ToEntity() *entity.CardMetrics {
	return &entity.CardMetrics{
		Id:       c.Id,
		CardId:   c.CardId,
		Likes:    c.Likes,
		Dislikes: c.Dislikes,
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

type CollectionUserProgress struct {
	Id           uuid.UUID `gorm:"primary_key;column:id"`
	CollectionId uuid.UUID `gorm:"column:collection_id"`
	UserId       uuid.UUID `gorm:"column:user_id"`
	Mastered     uint32    `gorm:"column:mastered"`
	Reviewing    uint32    `gorm:"column:reviewing"`
	Learning     uint32    `gorm:"column:learning"`
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
