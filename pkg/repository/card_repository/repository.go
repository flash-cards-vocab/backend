package card_repository

import (
	"errors"
	"math"
	"time"

	repositoryIntf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db        *gorm.DB
	tableName string
}

func New(db *gorm.DB) repositoryIntf.CardRepository {
	return &repository{db: db, tableName: "card"}
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
	return r.db.Table(r.tableName).Create(data).Error
}

func (r *repository) CreateMultipleCards(collectionId uuid.UUID, cards []*entity.Card, userId uuid.UUID) error {
	tx := r.db.Begin()
	cardsModels := []*Card{}
	for _, card := range cards {
		cardsModels = append(cardsModels, &Card{
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

	err := r.db.Table(r.tableName).Create(cardsModels).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	cardUserProgress := []*CardUserProgress{}
	for _, card := range cardsModels {
		cardUserProgress = append(cardUserProgress, &CardUserProgress{
			Id:            uuid.New(),
			CardId:        card.Id,
			UserId:        userId,
			Status:        entity.CardUserProgressType_None,
			LearningCount: 0,
		})
	}
	err = r.db.Table("card_user_progress").Create(cardUserProgress).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	cardMetrics := []*CardMetrics{}
	for _, card := range cardsModels {
		cardMetrics = append(cardMetrics, &CardMetrics{
			Id:       uuid.New(),
			CardId:   card.Id,
			Likes:    0,
			Dislikes: 0,
		})
	}
	err = r.db.Table("card_metrics").Create(cardMetrics).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	collectionCards := []*CollectionCards{}
	for _, card := range cardsModels {
		collectionCards = append(collectionCards, &CollectionCards{
			Id:           uuid.New(),
			CardId:       card.Id,
			CollectionId: collectionId,
		})
	}
	err = r.db.Table("collection_cards").Create(collectionCards).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *repository) AssignCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error {
	collectionCard := &CollectionCards{
		CollectionId: collectionId,
		CardId:       cardId,
	}
	return r.db.Table("collection_cards").Create(collectionCard).Error
}

func (r *repository) KnowCard(collectionId, cardId, userId uuid.UUID) error {
	tx := r.db.Begin()
	collectionUserProgress := CollectionUserProgress{}
	err := r.db.
		Table("collection_user_progress").
		Where("collection_id=? AND user_id=?", collectionId, userId).
		First(&collectionUserProgress).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collectionUserProgress = CollectionUserProgress{
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
				Create(&collectionUserProgress).
				Error
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	collectionCard := CardUserProgress{}
	err = r.db.
		Table("card_user_progress").
		Where("card_id=? AND user_id=?", cardId, userId).
		First(&collectionCard).
		Error
	if err != nil {
		// if user had no interactions with this card, create one
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collectionCard = CardUserProgress{
				Id:            uuid.New(),
				CardId:        cardId,
				UserId:        userId,
				Status:        entity.CardUserProgressType_Mastered,
				LearningCount: 3,
			}
			err := r.db.
				Table("card_user_progress").
				Create(&collectionCard).
				Error
			if err != nil {
				tx.Rollback()
				return err
			}

			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"mastered": collectionUserProgress.Mastered + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else {
			return err
		}
	} else {
		if collectionCard.LearningCount == 2 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Mastered,
					"learning_count": collectionCard.LearningCount + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"mastered":  collectionUserProgress.Mastered + 1,
					"reviewing": math.Max(float64(collectionUserProgress.Reviewing), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if collectionCard.LearningCount == 1 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Reviewing,
					"learning_count": collectionCard.LearningCount + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			// err = r.db.Table("collection_user_progress").
			// 	Where("collection_id=? AND user_id=?", collectionId, userId).
			// 	Updates(map[string]interface{}{
			// 		"reviewing": collectionUserProgress.Reviewing + 1,
			// 		"learning":  collectionUserProgress.Mastered + 1,
			// 	}).Error
			// if err != nil {
			// 	tx.Rollback()
			// 	return err
			// }

		} else if collectionCard.LearningCount == 0 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Reviewing,
					"learning_count": collectionCard.LearningCount + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"learning":  math.Max(float64(collectionUserProgress.Learning), 1) - 1,
					"reviewing": collectionUserProgress.Reviewing + 1,
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
	collectionUserProgress := CollectionUserProgress{}
	err := r.db.
		Table("collection_user_progress").
		Where("collection_id=? AND user_id=?", collectionId, userId).
		First(&collectionUserProgress).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collectionUserProgress = CollectionUserProgress{
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
				Create(&collectionUserProgress).
				Error
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	collectionCard := CardUserProgress{}
	err = r.db.
		Table("card_user_progress").
		Where("card_id=? AND user_id=?", cardId, userId).
		First(&collectionCard).
		Error

	if err != nil {
		// if user had no interactions with this card, craete one
		if errors.Is(err, gorm.ErrRecordNotFound) {
			collectionCard = CardUserProgress{
				Id:            uuid.New(),
				CardId:        cardId,
				UserId:        userId,
				Status:        entity.CardUserProgressType_Learning,
				LearningCount: 0,
			}
			err := r.db.
				Table("card_user_progress").
				Create(&collectionCard).
				Error
			if err != nil {
				tx.Rollback()
				return err
			}

			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"learning": collectionUserProgress.Learning + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else {
			return err
		}
	} else {
		if collectionCard.LearningCount == 3 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Reviewing,
					"learning_count": math.Max(float64(collectionCard.LearningCount), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"mastered":  math.Max(float64(collectionUserProgress.Mastered), 1) - 1,
					"reviewing": collectionUserProgress.Reviewing + 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if collectionCard.LearningCount == 2 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"learning_count": math.Max(float64(collectionCard.LearningCount), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if collectionCard.LearningCount == 1 {
			err = r.db.Table("card_user_progress").
				Where("card_id=? AND user_id=?", cardId, userId).
				Updates(map[string]interface{}{
					"status":         entity.CardUserProgressType_Learning,
					"learning_count": math.Max(float64(collectionCard.LearningCount), 1) - 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			err = r.db.Table("collection_user_progress").
				Where("collection_id=? AND user_id=?", collectionId, userId).
				Updates(map[string]interface{}{
					"learning":  collectionUserProgress.Learning + 1,
					"reviewing": math.Max(float64(collectionUserProgress.Reviewing), 1) - 1,
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
