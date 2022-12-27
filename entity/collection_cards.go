package entity

import (
	"github.com/google/uuid"
)

type CollectionCards struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CardId       uuid.UUID `json:"cardId,omitempty"`
	CollectionId uuid.UUID `json:"collectionId,omitempty"`
}
