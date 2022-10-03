package collection_usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type usecase struct {
	collection_repo repository.CollectionRepository
	card_repo       repository.CardRepository
	user_repo       repository.UserRepository
}

func New(
	collection_repo repository.CollectionRepository,
	card_repo repository.CardRepository,
	user_repo repository.UserRepository) UseCase {
	return &usecase{
		collection_repo: collection_repo,
		card_repo:       card_repo,
		user_repo:       user_repo,
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
		collection_metrics, err := uc.collection_repo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_user_progress, err := uc.collection_repo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collection_user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		totalCards, err := uc.collection_repo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collection_reponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       "You",
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collection_user_metrics.Starred,
			Likes:            collection_metrics.Likes,
			Dislikes:         collection_metrics.Dislikes,
			Views:            collection_metrics.Views,
			Mastered:         collection_user_progress.Mastered,
			Reviewing:        collection_user_progress.Reviewing,
			Learning:         collection_user_progress.Learning,
			IsLikedByUser:    collection_user_metrics.Liked,
			IsDislikedByUser: collection_user_metrics.Disliked,
			IsViewedByUser:   collection_user_metrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collection_reponses = append(collection_reponses, collection_reponse)
	}

	return collection_reponses, err
}

func (uc *usecase) GetCollectionFullUserMetrics(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error) {
	collection_metrics, err := uc.collection_repo.GetCollectionMetrics(id)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	collection_user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	collection_reponse := &entity.CollectionFullUserMetricsResponse{
		CollectionId: id,
		Likes:        collection_metrics.Likes,
		Dislikes:     collection_metrics.Dislikes,
		Views:        collection_metrics.Views,
		UserId:       userId,
		Liked:        collection_user_metrics.Liked,
		Disliked:     collection_user_metrics.Disliked,
		Viewed:       collection_user_metrics.Viewed,
		Starred:      collection_user_metrics.Starred,
	}
	return collection_reponse, nil
}

func (uc *usecase) GetRecommendedCollectionsPreview(userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {
	collection_reponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collection_repo.GetRecommendedCollectionsPreview(userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collection_author, err := uc.user_repo.GetUserById(collection.AuthorId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_metrics, err := uc.collection_repo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collection_user_progress, err := uc.collection_repo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		totalCards, err := uc.collection_repo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collection_reponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       collection_author.Name,
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collection_user_metrics.Starred,
			Likes:            collection_metrics.Likes,
			Dislikes:         collection_metrics.Dislikes,
			Views:            collection_metrics.Views,
			Mastered:         collection_user_progress.Mastered,
			Reviewing:        collection_user_progress.Reviewing,
			Learning:         collection_user_progress.Learning,
			IsLikedByUser:    collection_user_metrics.Liked,
			IsDislikedByUser: collection_user_metrics.Disliked,
			IsViewedByUser:   collection_user_metrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collection_reponses = append(collection_reponses, collection_reponse)
	}

	return collection_reponses, err
}

func (uc *usecase) GetLikedCollectionsPreview(userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {
	collection_reponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collection_repo.GetLikedCollectionsPreview(userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collection_author, err := uc.user_repo.GetUserById(collection.AuthorId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_metrics, err := uc.collection_repo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collection_user_progress, err := uc.collection_repo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		totalCards, err := uc.collection_repo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collection_reponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       collection_author.Name,
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collection_user_metrics.Starred,
			Likes:            collection_metrics.Likes,
			Dislikes:         collection_metrics.Dislikes,
			Views:            collection_metrics.Views,
			Mastered:         collection_user_progress.Mastered,
			Reviewing:        collection_user_progress.Reviewing,
			Learning:         collection_user_progress.Learning,
			IsLikedByUser:    collection_user_metrics.Liked,
			IsDislikedByUser: collection_user_metrics.Disliked,
			IsViewedByUser:   collection_user_metrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collection_reponses = append(collection_reponses, collection_reponse)
	}

	return collection_reponses, err
}

func (uc *usecase) GetStarredCollectionsPreview(userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {
	collection_reponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collection_repo.GetStarredCollectionsPreview(userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collection_author, err := uc.user_repo.GetUserById(collection.AuthorId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_metrics, err := uc.collection_repo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collection_user_progress, err := uc.collection_repo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		totalCards, err := uc.collection_repo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collection_reponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       collection_author.Name,
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collection_user_metrics.Starred,
			Likes:            collection_metrics.Likes,
			Dislikes:         collection_metrics.Dislikes,
			Views:            collection_metrics.Views,
			Mastered:         collection_user_progress.Mastered,
			Reviewing:        collection_user_progress.Reviewing,
			Learning:         collection_user_progress.Learning,
			IsLikedByUser:    collection_user_metrics.Liked,
			IsDislikedByUser: collection_user_metrics.Disliked,
			IsViewedByUser:   collection_user_metrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collection_reponses = append(collection_reponses, collection_reponse)
	}

	return collection_reponses, err
}

func (uc *usecase) GetCollectionWithCards(id, userId uuid.UUID, page, size int) (*entity.GetCollectionWithCardsResponse, error) {
	var err error
	collection, err := uc.collection_repo.GetCollection(id)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	collection_progress, err := uc.collection_repo.GetCollectionUserProgress(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	limit := size
	offset := (page - 1) * size
	cards, err := uc.collection_repo.GetCollectionCards(id, limit, offset)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	collection_reponses := &entity.GetCollectionWithCardsResponse{
		Id:         collection.Id,
		Name:       collection.Name,
		Mastered:   collection_progress.Mastered,
		Reviewing:  collection_progress.Reviewing,
		Learning:   collection_progress.Learning,
		TotalCards: cards.Total,
		Cards:      cards.Cards,
	}
	return collection_reponses, nil
}

func (uc *usecase) StarCollectionById(id, userId uuid.UUID) error {
	err := uc.collection_repo.StarCollectionById(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	return nil
}

func (uc *usecase) LikeCollectionById(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error) {
	_, err := uc.collection_repo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			err = uc.collection_repo.CreateCollectionUserMetrics(id, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	isLiked, isDisliked, err := uc.collection_repo.IsCollectionLikedOrDislikedByUser(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	if isDisliked {
		err := uc.collection_repo.CollectionDislikeInteraction(id, userId, true)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	err = uc.collection_repo.CollectionLikeInteraction(id, userId, isLiked)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	metrics, err := uc.collection_repo.GetCollectionMetrics(id)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	return &entity.CollectionFullUserMetricsResponse{
		CollectionId: id,
		Likes:        metrics.Likes,
		Dislikes:     metrics.Dislikes,
		Views:        metrics.Views,
		UserId:       userId,
		Liked:        user_metrics.Liked,
		Disliked:     user_metrics.Disliked,
		Viewed:       user_metrics.Viewed,
		Starred:      user_metrics.Starred,
	}, nil
}

func (uc *usecase) DislikeCollectionById(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error) {
	_, err := uc.collection_repo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			err = uc.collection_repo.CreateCollectionUserMetrics(id, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	isLiked, isDisliked, err := uc.collection_repo.IsCollectionLikedOrDislikedByUser(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	if isLiked {
		err := uc.collection_repo.CollectionLikeInteraction(id, userId, true)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	err = uc.collection_repo.CollectionDislikeInteraction(id, userId, isDisliked)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	metrics, err := uc.collection_repo.GetCollectionMetrics(id)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	return &entity.CollectionFullUserMetricsResponse{
		CollectionId: id,
		Likes:        metrics.Likes,
		Dislikes:     metrics.Dislikes,
		Views:        metrics.Views,
		UserId:       userId,
		Liked:        user_metrics.Liked,
		Disliked:     user_metrics.Disliked,
		Viewed:       user_metrics.Viewed,
		Starred:      user_metrics.Starred,
	}, nil
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

func (uc *usecase) SearchCollectionByName(text string, userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {

	collection_reponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collection_repo.SearchCollectionByName(text, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collection_author, err := uc.user_repo.GetUserById(collection.AuthorId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_metrics, err := uc.collection_repo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collection_user_progress, err := uc.collection_repo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collection_user_metrics, err := uc.collection_repo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		totalCards, err := uc.collection_repo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repository.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collection_reponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       collection_author.Name,
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collection_user_metrics.Starred,
			Likes:            collection_metrics.Likes,
			Dislikes:         collection_metrics.Dislikes,
			Views:            collection_metrics.Views,
			Mastered:         collection_user_progress.Mastered,
			Reviewing:        collection_user_progress.Reviewing,
			Learning:         collection_user_progress.Learning,
			IsLikedByUser:    collection_user_metrics.Liked,
			IsDislikedByUser: collection_user_metrics.Disliked,
			IsViewedByUser:   collection_user_metrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collection_reponses = append(collection_reponses, collection_reponse)
	}

	return collection_reponses, err

}

func (uc *usecase) CreateCollection(collection entity.Collection, cards []*entity.Card, userId uuid.UUID) error {
	createdCollection, err := uc.collection_repo.CreateCollection(collection)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	err = uc.card_repo.CreateMultipleCards(createdCollection.Id, cards, userId)
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
