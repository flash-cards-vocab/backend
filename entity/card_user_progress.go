package entity

import (
	"github.com/google/uuid"
)

type CardUserProgressType string

const (
	CardUserProgressType_Mastered  CardUserProgressType = "mastered"
	CardUserProgressType_Reviewing CardUserProgressType = "reviewing"
	CardUserProgressType_Learning  CardUserProgressType = "learning"
	CardUserProgressType_None      CardUserProgressType = "none"
)

type CardUserProgress struct {
	Id     uuid.UUID            `json:"id,omitempty"`
	CardId uuid.UUID            `json:"cardId,omitempty"`
	UserId uuid.UUID            `json:"userId,omitempty"`
	Status CardUserProgressType `json:"learning,omitempty"`
}
