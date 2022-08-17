package entity

import (
	"github.com/google/uuid"
)

type CollectionCards struct {
	Id           uuid.UUID `json:"id,omitempty"`
	CardId       uuid.UUID `json:"card_id,omitempty"`
	CollectionId uuid.UUID `json:"collection_id,omitempty"`
}
