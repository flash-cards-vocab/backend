package collection_repository

import (
	"time"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Collection struct {
	Id        uuid.UUID      `gorm:"primary_key;column:id"`
	Name      string         `gorm:"column:name"`
	AuthorId  uuid.UUID      `gorm:"column:author_id"`
	Topics    pq.StringArray `gorm:"type:text[];column:topics"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt *time.Time     `gorm:"column:deleted_at"`
}

func (c *Collection) ToEntity() *entity.Collection {
	return &entity.Collection{
		Id:        c.Id,
		Name:      c.Name,
		Topics:    c.Topics,
		AuthorId:  c.AuthorId,
		CreatedAt: c.CreatedAt,
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

type CollectionUserMetrics struct {
	Id           uuid.UUID `gorm:"primary_key;column:id"`
	UserId       uuid.UUID `gorm:"column:user_id"`
	CollectionId uuid.UUID `gorm:"column:collection_id"`
	Liked        bool      `gorm:"column:liked"`
	Disliked     bool      `gorm:"column:disliked"`
	Viewed       bool      `gorm:"column:viewed"`
	Starred      bool      `gorm:"column:starred"`
}

func (c *CollectionUserMetrics) ToEntity() *entity.CollectionUserMetrics {
	return &entity.CollectionUserMetrics{
		Id:           c.Id,
		UserId:       c.UserId,
		CollectionId: c.CollectionId,
		Liked:        c.Liked,
		Disliked:     c.Disliked,
		Viewed:       c.Viewed,
		Starred:      c.Starred,
	}
}

type CollectionMetrics struct {
	Id           uuid.UUID `gorm:"primary_key;column:id"`
	CollectionId uuid.UUID `gorm:"column:collection_id"`
	Likes        uint32    `gorm:"column:likes"`
	Dislikes     uint32    `gorm:"column:dislikes"`
	Views        uint32    `gorm:"column:views"`
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
	Id           uuid.UUID `gorm:"primary_key;column:id"`
	CardId       uuid.UUID `gorm:"column:card_id"`
	CollectionId uuid.UUID `gorm:"column:collection_id"`
}

func (c *CollectionCards) ToEntity() *entity.CollectionCards {
	return &entity.CollectionCards{
		Id:           c.Id,
		CardId:       c.CardId,
		CollectionId: c.CollectionId,
	}
}

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

func (c Card) ToArrayEntity(cards []*Card) []*entity.Card {
	res := []*entity.Card{}
	for _, card := range cards {
		res = append(res, &entity.Card{
			Id:         card.Id,
			Word:       card.Word,
			ImageUrl:   card.ImageUrl,
			Definition: card.Definition,
			Sentence:   card.Sentence,
			Antonyms:   card.Antonyms,
			Synonyms:   card.Synonyms,
		})
	}
	return res
}
