package collection_usecase

import (
	"errors"
	"fmt"

	"github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type usecase struct {
	collection_repo repository.CollectionRepository
	card_repo       repository.CardRepository
}

func New(collection_repo repository.CollectionRepository, card_repo repository.CardRepository) UseCase {
	return &usecase{
		collection_repo: collection_repo,
		card_repo:       card_repo,
	}
}

func (uc *usecase) GetMyCollections(userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {
	collection_reponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collection_repo.GetMyCollections(userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		// collection_metrics, err := uc.collection_repo.GetCollectionMetrics(collection.Id, userId)
		// if err != nil {
		// 	if errors.Is(err, repository.ErrCollectionNotFound) {
		// 		return nil, ErrNotFound
		// 	}
		// 	logrus.Errorf("%w: %v", ErrUnexpected, err)
		// 	return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		// }

		// collection_user_progress, err := uc.collection_repo.GetCollectionUserProgress(collection.Id, userId)
		// if err != nil {
		// 	if errors.Is(err, repository.ErrCollectionNotFound) {
		// 		return nil, ErrNotFound
		// 	}
		// 	logrus.Errorf("%w: %v", ErrUnexpected, err)
		// 	return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		// }
		// collection_user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(collection.Id, userId)
		// if err != nil {
		// 	if errors.Is(err, repository.ErrCollectionNotFound) {
		// 		return nil, ErrNotFound
		// 	}
		// 	logrus.Errorf("%w: %v", ErrUnexpected, err)
		// 	return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		// }

		collection_reponse := &entity.UserCollectionResponse{
			Id:     collection.Id,
			Name:   collection.Name,
			Topics: collection.Topics,
			// Likes:            collection_metrics.Likes,
			// Dislikes:         collection_metrics.Dislikes,
			// Views:            collection_metrics.Views,
			// Mastered:         collection_user_progress.Mastered,
			// Reviewing:        collection_user_progress.Reviewing,
			// Learning:         collection_user_progress.Learning,
			// IsLikedByUser:    collection_user_metrics.Liked,
			// IsDislikedByUser: collection_user_metrics.Disliked,
			// IsViewedByUser:   collection_user_metrics.Viewed,
		}

		collection_reponses = append(collection_reponses, collection_reponse)
	}

	return collection_reponses, err
}

func (uc *usecase) LikeCollectionById(id, userId uuid.UUID) error {
	isLiked, err := uc.collection_repo.IsCollectionLikedByUser(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	isDisliked, err := uc.collection_repo.IsCollectionDislikedByUser(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	if isDisliked {
		err := uc.collection_repo.CollectionDislikeInteraction(id, userId, false)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	err = uc.collection_repo.CollectionLikeInteraction(id, userId, isLiked)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	return nil
}

func (uc *usecase) DislikeCollectionById(id, userId uuid.UUID) error {
	isDisliked, err := uc.collection_repo.IsCollectionDislikedByUser(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	isLiked, err := uc.collection_repo.IsCollectionLikedByUser(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	if isLiked {
		err := uc.collection_repo.CollectionLikeInteraction(id, userId, false)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	err = uc.collection_repo.CollectionDislikeInteraction(id, userId, isDisliked)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	return nil
}

func (uc *usecase) ViewCollectionById(id, userId uuid.UUID) error {
	isViewed, err := uc.collection_repo.IsCollectionViewedByUser(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	if !isViewed {
		err := uc.collection_repo.ViewCollection(id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}
	return nil
}

func (uc *usecase) SearchCollectionByName(text string) ([]*entity.Collection, error) {
	panic("Not implemented")
}

func (uc *usecase) CreateCollection(collection entity.Collection, cards []*entity.Card) error {
	createdCollection, err := uc.collection_repo.CreateCollection(collection)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	err = uc.card_repo.CreateMultipleCards(createdCollection.Id, cards)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	return nil
}

func (uc *usecase) UpdateCollectionUserProgress(id uuid.UUID, mastered, reviewing, learning uint32) error {
	panic("Not implemented")
}
