package repository

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

var ErrCardNotFound = errors.New("Card not found")

type CardRepository interface {
	CreateSingleCard(card entity.Card) error
	CreateMultipleCards(collectionId uuid.UUID, card []*entity.Card, userId uuid.UUID) error
	RemoveMultipleCardsFromCollection(cardsToRemove []*entity.CollectionCards) error
	AssignCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error
	KnowCard(collectionId, cardId, userId uuid.UUID) error
	DontKnowCard(collectionId, cardId, userId uuid.UUID) error
	GetUserCardsStatistics(userId uuid.UUID) (*entity.UserCardStatistics, error)
	GetCardsByWord(word string, limit, offset int) ([]*entity.Card, error)
}
