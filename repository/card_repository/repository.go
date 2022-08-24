package card_repository

import (
	"time"

	repository_intf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db         *gorm.DB
	table_name string
}

func New(db *gorm.DB) repository_intf.CardRepository {
	return &repository{db: db, table_name: "card"}
}

func (r *repository) CreateSingleCard(card entity.Card) error {
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

func (r *repository) CreateMultipleCards(collectionId uuid.UUID, cards []*entity.Card) error {

	cards_models := []*Card{}
	for _, card := range cards {
		cards_models = append(cards_models, &Card{
			Word:       card.Word,
			ImageUrl:   card.ImageUrl,
			Definition: card.Definition,
			Sentence:   card.Sentence,
			Antonyms:   card.Antonyms,
			Synonyms:   card.Synonyms,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	err := r.db.Table(r.table_name).Create(cards_models).Error
	if err != nil {
		return err
	}

	collection_cards := []*CollectionCards{}
	for _, card := range cards_models {
		collection_cards = append(collection_cards, &CollectionCards{
			CardId:       card.Id,
			CollectionId: collectionId,
		})
	}
	return r.db.Table("collection_cards").Create(collection_cards).Error
}

func (r *repository) AssignCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error {
	collection_card := &CollectionCards{
		CollectionId: collectionId,
		CardId:       cardId,
	}
	return r.db.Table("collection_cards").Create(collection_card).Error
}
