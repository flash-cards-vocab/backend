package card_repository

import (
	"time"

	repository_intf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"gorm.io/gorm"
)

type repository struct {
	db         *gorm.DB
	table_name string
}

func New(db *gorm.DB) repository_intf.CardRepository {
	return &repository{db: db, table_name: "card"}
}

func (r *repository) CreateCard(card entity.Card) error {
	data := &Card{
		Word:       card.Word,
		ImageUrl:   card.ImageUrl,
		Definition: card.Definition,
		Sentence:   card.Sentence,
		Antonyms:   card.Antonyms,
		Synonyms:   card.Synonyms,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return r.db.Table(r.table_name).Create(data).Error
}
