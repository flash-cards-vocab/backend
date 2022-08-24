package repository

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

var ErrCardNotFound = errors.New("Card not found")

type CardRepository interface {
	CreateSingleCard(card entity.Card) error
	CreateMultipleCards(collectionId uuid.UUID, card []*entity.Card) error
	AssignCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error
}
