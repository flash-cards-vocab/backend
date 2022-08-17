package repository

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
)

var ErrCardNotFound = errors.New("Card not found")

type CardRepository interface {
	CreateCard(card entity.Card) error
}
