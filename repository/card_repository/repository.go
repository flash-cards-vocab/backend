package card_repository

import (
	"errors"
	"math"
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
		Id:         uuid.New(),
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

func (r *repository) CreateMultipleCards(collectionId uuid.UUID, cards []*entity.Card, userId uuid.UUID) error {
	tx := r.db.Begin()
	cards_models := []*Card{}
	for _, card := range cards {
		cards_models = append(cards_models, &Card{
			Id:         uuid.New(),
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
		tx.Rollback()
		return err
	}

	card_user_progress := []*CardUserProgress{}
	for _, card := range cards_models {
		card_user_progress = append(card_user_progress, &CardUserProgress{
			Id:            uuid.New(),
			CardId:        card.Id,
			UserId:        userId,
			Status:        entity.CardUserProgressType_None,
			LearningCount: 0,
		})
	}
	err = r.db.Table("card_user_progress").Create(card_user_progress).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	card_metrics := []*CardMetrics{}
	for _, card := range cards_models {
		card_metrics = append(card_metrics, &CardMetrics{
			Id:       uuid.New(),
			CardId:   card.Id,
			Likes:    0,
			Dislikes: 0,
		})
	}
	err = r.db.Table("card_metrics").Create(card_metrics).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	collection_cards := []*CollectionCards{}
	for _, card := range cards_models {
		collection_cards = append(collection_cards, &CollectionCards{
			Id:           uuid.New(),
			CardId:       card.Id,
			CollectionId: collectionId,
		})
	}
	err = r.db.Table("collection_cards").Create(collection_cards).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *repository) AssignCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error {
	collection_card := &CollectionCards{
		CollectionId: collectionId,
		CardId:       cardId,
	}
	return r.db.Table("collection_cards").Create(collection_card).Error
}

func (r *repository) KnowCard(collectionId, cardId, userId uuid.UUID) error {
	tx := r.db.Begin()
	collection_user_progress := CollectionUserProgress{}
	err := r.db.
		Table("collection_user_progress").
		Where("collection_id=? AND user_id=?", collectionId, userId).
		First(&collection_user_progress).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collection_user_progress = CollectionUserProgress{
				Id:           uuid.New(),
				CollectionId: collectionId,
				UserId:       userId,
				Mastered:     0,
				Reviewing:    0,
				Learning:     0,
			}
			err = r.db.
				Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Create(&collection_user_progress).
				Error
		} else {
			return err
		}
	}

	collection_card := CardUserProgress{}
	err = r.db.
		Table("card_user_progress").
		Where("card_id=? AND user_id=?", cardId, userId).
		First(&collection_card).
		Error
	if err != nil {
		// if user had no interactions with this card, create one
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collection_card = CardUserProgress{
				Id:            uuid.New(),
				CardId:        cardId,
				UserId:        userId,
				Status:        entity.CardUserProgressType_Mastered,
				LearningCount: 3,
			}
			err := r.db.
				Table("card_user_progress").
				Create(&collection_card).
				Error
			if err != nil {
				tx.Rollback()
				return err
			}

			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"mastered": collection_user_progress.Mastered + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else {
			return err
		}
	} else {
		if collection_card.LearningCount == 2 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Mastered,
					"learning_count": collection_card.LearningCount + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"mastered":  collection_user_progress.Mastered + 1,
					"reviewing": math.Max(float64(collection_user_progress.Reviewing), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if collection_card.LearningCount == 1 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Learning,
					"learning_count": collection_card.LearningCount + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		} else if collection_card.LearningCount == 0 {
			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"learning":  math.Max(float64(collection_user_progress.Learning), 1) - 1,
					"reviewing": collection_user_progress.Reviewing + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil

	//
}

func (r *repository) DontKnowCard(collectionId, cardId, userId uuid.UUID) error {
	tx := r.db.Begin()
	collection_user_progress := CollectionUserProgress{}

	err := r.db.
		Table("collection_user_progress").
		Where("collection_id=? AND user_id=?", collectionId, userId).
		First(&collection_user_progress).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collection_user_progress = CollectionUserProgress{
				Id:           uuid.New(),
				CollectionId: collectionId,
				UserId:       userId,
				Mastered:     0,
				Reviewing:    0,
				Learning:     0,
			}
			err = r.db.
				Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Create(&collection_user_progress).
				Error
		} else {
			return err
		}
	}
	collection_card := CardUserProgress{}
	err = r.db.
		Table("card_user_progress").
		Where("card_id=? AND user_id=?", cardId, userId).
		Scan(&collection_card).
		Error

	if err != nil {
		// if user had no interactions with this card, craete one
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collection_card = CardUserProgress{
				Id:            uuid.New(),
				CardId:        cardId,
				UserId:        userId,
				Status:        entity.CardUserProgressType_Learning,
				LearningCount: 0,
			}
			err := r.db.
				Table("card_user_progress").
				Create(&collection_card).
				Error
			if err != nil {
				tx.Rollback()
				return err
			}

			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"learning": collection_user_progress.Learning + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else {
			return err
		}
	} else {
		if collection_card.LearningCount == 3 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Reviewing,
					"learning_count": math.Max(float64(collection_card.LearningCount), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"mastered":  math.Max(float64(collection_user_progress.Mastered), 1) - 1,
					"reviewing": collection_user_progress.Reviewing + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if collection_card.LearningCount == 2 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"learning_count": math.Max(float64(collection_card.LearningCount), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if collection_card.LearningCount == 1 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Learning,
					"learning_count": math.Max(float64(collection_card.LearningCount), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"learning":  collection_user_progress.Learning + 1,
					"reviewing": math.Max(float64(collection_user_progress.Reviewing), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	}
	tx.Commit()
	return nil
}
