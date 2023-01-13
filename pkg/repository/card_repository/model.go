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

type CardWithOccurence struct {
	Id         uuid.UUID `gorm:"primary_key;column:id"`
	Word       string    `gorm:"column:word"`
	ImageUrl   string    `gorm:"column:imageUrl"`
	Definition string    `gorm:"column:definition"`
	Sentence   string    `gorm:"column:sentence"`
	Antonyms   string    `gorm:"column:antonyms"`
	Synonyms   string    `gorm:"column:synonyms"`
	AuthorId   uuid.UUID `gorm:"column:authorId"`
	Occurence  int       `gorm:"column:occurence"`
}

func (c *CardWithOccurence) ToEntity() *entity.CardWithOccurence {
	return &entity.CardWithOccurence{
		Id:         c.Id,
		Word:       c.Word,
		ImageUrl:   c.ImageUrl,
		Definition: c.Definition,
		Sentence:   c.Sentence,
		Antonyms:   c.Antonyms,
		Synonyms:   c.Synonyms,
		AuthorId:   c.AuthorId,
		Occurence:  c.Occurence,
	}
}

func (c CardWithOccurence) ToArrayEntity(cards []*CardWithOccurence) []*entity.CardWithOccurence {
	res := []*entity.CardWithOccurence{}
	for _, card := range cards {
		res = append(res, &entity.CardWithOccurence{
			Id:         card.Id,
			Word:       card.Word,
			ImageUrl:   card.ImageUrl,
			Definition: card.Definition,
			Sentence:   card.Sentence,
			Antonyms:   card.Antonyms,
			Synonyms:   card.Synonyms,
			AuthorId:   card.AuthorId,
			Occurence:  card.Occurence,
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
func (c CollectionCards) FromArrayEntity(cards []*entity.CollectionCards) []*CollectionCards {
	res := []*CollectionCards{}
	for _, card := range cards {
		res = append(res, &CollectionCards{
			CardId:       card.CardId,
			CollectionId: card.CollectionId,
		})
	}
	return res
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

type UserCardsStatistics struct {
	CardsCreated uint32 `gorm:"column:cards_created"`
	Mastered     uint32 `gorm:"column:mastered"`
	Reviewing    uint32 `gorm:"column:reviewing"`
	Learning     uint32 `gorm:"column:learning"`
}
