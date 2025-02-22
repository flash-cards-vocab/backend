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
	Id           uuid.UUID  `gorm:"primary_key;column:id"`
	CollectionId uuid.UUID  `gorm:"column:collection_id"`
	UserId       uuid.UUID  `gorm:"column:user_id"`
	Mastered     uint32     `gorm:"column:mastered"`
	Reviewing    uint32     `gorm:"column:reviewing"`
	Learning     uint32     `gorm:"column:learning"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
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
	Id           uuid.UUID  `gorm:"primary_key;column:id"`
	UserId       uuid.UUID  `gorm:"column:user_id"`
	CollectionId uuid.UUID  `gorm:"column:collection_id"`
	Liked        bool       `gorm:"column:liked"`
	Disliked     bool       `gorm:"column:disliked"`
	Viewed       bool       `gorm:"column:viewed"`
	Starred      bool       `gorm:"column:starred"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
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
	Id           uuid.UUID  `gorm:"primary_key;column:id"`
	CollectionId uuid.UUID  `gorm:"column:collection_id"`
	Likes        uint32     `gorm:"column:likes"`
	Dislikes     uint32     `gorm:"column:dislikes"`
	Views        uint32     `gorm:"column:views"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
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
	Id           uuid.UUID  `gorm:"primary_key;column:id"`
	CardId       uuid.UUID  `gorm:"column:card_id"`
	CollectionId uuid.UUID  `gorm:"column:collection_id"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
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
	AuthorId   uuid.UUID  `gorm:"column:author_id"`
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
		AuthorId:   c.AuthorId,
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
			AuthorId:   card.AuthorId,
		})
	}
	return res
}

type CardForUser struct {
	Id         uuid.UUID `gorm:"column:id"`
	Word       string    `gorm:"column:word"`
	ImageUrl   string    `gorm:"column:image_url"`
	Definition string    `gorm:"column:definition"`
	Sentence   string    `gorm:"column:sentence"`
	Antonyms   string    `gorm:"column:antonyms"`
	Synonyms   string    `gorm:"column:synonyms"`
	Status     string    `gorm:"column:status"`
	AuthorId   uuid.UUID `gorm:"column:author_id"`
}

func (c *CardForUser) ToEntity() *entity.CardForUser {
	return &entity.CardForUser{
		Id:         c.Id,
		Word:       c.Word,
		ImageUrl:   c.ImageUrl,
		Definition: c.Definition,
		Sentence:   c.Sentence,
		Antonyms:   c.Antonyms,
		Synonyms:   c.Synonyms,
		Status:     c.Status,
		AuthorId:   c.AuthorId,
	}
}

func (c CardForUser) ToArrayEntity(cards []*CardForUser) []*entity.CardForUser {
	res := []*entity.CardForUser{}
	for _, card := range cards {
		res = append(res, &entity.CardForUser{
			Id:         card.Id,
			Word:       card.Word,
			ImageUrl:   card.ImageUrl,
			Definition: card.Definition,
			Sentence:   card.Sentence,
			Antonyms:   card.Antonyms,
			Synonyms:   card.Synonyms,
			Status:     card.Status,
			AuthorId:   c.AuthorId,
		})
	}
	return res
}

type CardUserProgress struct {
	Id            uuid.UUID                   `gorm:"primary_key;column:id"`
	CardId        uuid.UUID                   `gorm:"column:card_id"`
	UserId        uuid.UUID                   `gorm:"column:user_id"`
	Status        entity.CardUserProgressType `gorm:"column:status"`
	LearningCount uint32                      `gorm:"column:learning_count"`
	CreatedAt     time.Time                   `gorm:"column:created_at"`
	UpdatedAt     time.Time                   `gorm:"column:updated_at"`
	DeletedAt     *time.Time                  `gorm:"column:deleted_at"`
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
	Id        uuid.UUID  `gorm:"primary_key;column:id"`
	CardId    uuid.UUID  `gorm:"column:card_id"`
	Likes     uint32     `gorm:"column:likes"`
	Dislikes  uint32     `gorm:"column:dislikes"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (c *CardMetrics) ToEntity() *entity.CardMetrics {
	return &entity.CardMetrics{
		Id:       c.Id,
		CardId:   c.CardId,
		Likes:    c.Likes,
		Dislikes: c.Dislikes,
	}
}
