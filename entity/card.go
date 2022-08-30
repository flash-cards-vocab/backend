package entity

import (
	"github.com/google/uuid"
)

type Card struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Word       string    `json:"word,omitempty"`
	ImageUrl   string    `json:"image_url,omitempty"`
	Definition string    `json:"definition,omitempty"`
	Sentence   string    `json:"sentence,omitempty"`
	Antonyms   string    `json:"antonyms,omitempty"`
	Synonyms   string    `json:"synonyms,omitempty"`
}

type CardsPagination struct {
	Cards []*Card `json:"cards,omitempty"`
	Page  int     `json:"page,omitempty"`
	Size  int     `json:"size,omitempty"`
	Total int     `json:"total,omitempty"`
}
