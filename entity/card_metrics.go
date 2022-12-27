package entity

import (
	"github.com/google/uuid"
)

type CardMetrics struct {
	Id       uuid.UUID `json:"id,omitempty"`
	CardId   uuid.UUID `json:"cardId,omitempty"`
	Likes    uint32    `json:"likes,omitempty"`
	Dislikes uint32    `json:"dislikes,omitempty"`
}
