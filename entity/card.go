package entity

import (
	"github.com/google/uuid"
)

type Card struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Word       string    `json:"word,omitempty"`
	ImageUrl   string    `json:"imageUrl,omitempty"`
	Definition string    `json:"definition,omitempty"`
	Sentence   string    `json:"sentence,omitempty"`
	Antonyms   string    `json:"antonyms,omitempty"`
	Synonyms   string    `json:"synonyms,omitempty"`
	AuthorId   uuid.UUID `json:"authorId,omitempty"`
}
type CardWithOccurence struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Word       string    `json:"word,omitempty"`
	ImageUrl   string    `json:"imageUrl,omitempty"`
	Definition string    `json:"definition,omitempty"`
	Sentence   string    `json:"sentence,omitempty"`
	Antonyms   string    `json:"antonyms,omitempty"`
	Synonyms   string    `json:"synonyms,omitempty"`
	AuthorId   uuid.UUID `json:"authorId,omitempty"`
	Occurence  int       `json:"occurence,omitempty"`
}

type CardForUserPagination struct {
	CardForUser []*CardForUser `json:"cards,omitempty"`
	Page        int            `json:"page,omitempty"`
	Size        int            `json:"size,omitempty"`
	Total       int            `json:"total,omitempty"`
}

type CardForUser struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Word       string    `json:"word,omitempty"`
	ImageUrl   string    `json:"imageUrl,omitempty"`
	Definition string    `json:"definition,omitempty"`
	Sentence   string    `json:"sentence,omitempty"`
	Antonyms   string    `json:"antonyms,omitempty"`
	Synonyms   string    `json:"synonyms,omitempty"`
	Status     string    `json:"status,omitempty"`
	AuthorId   uuid.UUID `json:"authorId,omitempty"`
}

type CardUpdateType string

const (
	CardUpdateType_Create CardUpdateType = "Create"
	CardUpdateType_Update CardUpdateType = "Update"
	CardUpdateType_Remove CardUpdateType = "Remove"
)

type CardUpdate struct {
	Id         uuid.UUID      `json:"id,omitempty"`
	Word       string         `json:"word,omitempty"`
	ImageUrl   string         `json:"imageUrl,omitempty"`
	Definition string         `json:"definition,omitempty"`
	Sentence   string         `json:"sentence,omitempty"`
	Antonyms   string         `json:"antonyms,omitempty"`
	Synonyms   string         `json:"synonyms,omitempty"`
	Action     CardUpdateType `json:"action,omitempty"`
}

type CardSearch struct {
	UserCards   []*CardWithOccurence `json:"userCards,omitempty"`
	GlobalCards []*CardWithOccurence `json:"globalCards,omitempty"`
}
