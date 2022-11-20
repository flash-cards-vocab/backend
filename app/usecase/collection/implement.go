package collection_usecase

import (
	"errors"
	"fmt"
	"time"

	repositoryIntf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type usecase struct {
	collectionRepo repositoryIntf.CollectionRepository
	cardRepo       repositoryIntf.CardRepository
	userRepo       repositoryIntf.UserRepository
}

func New(
	collectionRepo repositoryIntf.CollectionRepository,
	cardRepo repositoryIntf.CardRepository,
	userRepo repositoryIntf.UserRepository) UseCase {
	return &usecase{
		collectionRepo: collectionRepo,
		cardRepo:       cardRepo,
		userRepo:       userRepo,
	}
}

func (uc *usecase) GetMyCollections(userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {
	collectionResponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collectionRepo.GetMyCollections(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositoryIntf.ErrCardNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collectionMetrics, err := uc.collectionRepo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collectionUserProgress, err := uc.collectionRepo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collectionUserMetrics, err := uc.collectionRepo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
				err = uc.collectionRepo.CreateCollectionUserMetrics(collection.Id, userId)
				if err != nil {
					logrus.Errorf("%w: %v", ErrUnexpected, err)
					return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
				}
			} else {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		}

		totalCards, err := uc.collectionRepo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collectionResponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       "You",
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collectionUserMetrics.Starred,
			Likes:            collectionMetrics.Likes,
			Dislikes:         collectionMetrics.Dislikes,
			Views:            collectionMetrics.Views,
			Mastered:         collectionUserProgress.Mastered,
			Reviewing:        collectionUserProgress.Reviewing,
			Learning:         collectionUserProgress.Learning,
			IsLikedByUser:    collectionUserMetrics.Liked,
			IsDislikedByUser: collectionUserMetrics.Disliked,
			IsViewedByUser:   collectionUserMetrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collectionResponses = append(collectionResponses, collectionResponse)
	}

	return collectionResponses, err
}

func (uc *usecase) GetCollectionUserProgress(id, userId uuid.UUID) (*entity.CollectionUserProgressResponse, error) {
	// collectionMetrics, err := uc.collectionRepo.GetCollectionMetrics(id)
	// if err != nil {
	// 	if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
	// 		return nil, ErrNotFound
	// 	}
	// 	logrus.Errorf("%w: %v", ErrUnexpected, err)
	// 	return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	// }

	// return collectionResponse, nil
	panic("Not implemented")
}

func (uc *usecase) GetCollectionFullUserMetrics(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error) {
	collectionMetrics, err := uc.collectionRepo.GetCollectionMetrics(id)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	collectionUserMetrics, err := uc.collectionRepo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
			err = uc.collectionRepo.CreateCollectionUserMetrics(id, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	collectionResponse := &entity.CollectionFullUserMetricsResponse{
		CollectionId: id,
		Likes:        collectionMetrics.Likes,
		Dislikes:     collectionMetrics.Dislikes,
		Views:        collectionMetrics.Views,
		UserId:       userId,
		Liked:        collectionUserMetrics.Liked,
		Disliked:     collectionUserMetrics.Disliked,
		Viewed:       collectionUserMetrics.Viewed,
		Starred:      collectionUserMetrics.Starred,
	}
	return collectionResponse, nil
}

func (uc *usecase) GetRecommendedCollectionsPreview(userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {
	collection_reponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collection_repo.GetRecommendedCollectionsPreview(userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collectionAuthor, err := uc.userRepo.GetUserById(collection.AuthorId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collectionMetrics, err := uc.collectionRepo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collectionUserProgress, err := uc.collectionRepo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionUserProgressNotFound) {
				err = uc.collectionRepo.CreateCollectionUserProgress(collection.Id, userId)
				if err != nil {
					logrus.Errorf("%w: %v", ErrUnexpected, err)
					return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
				}
			} else {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		}

		collectionUserMetrics, err := uc.collectionRepo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
				err = uc.collectionRepo.CreateCollectionUserMetrics(collection.Id, userId)
				if err != nil {
					logrus.Errorf("%w: %v", ErrUnexpected, err)
					return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
				}
			} else {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		}
		totalCards, err := uc.collectionRepo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collectionResponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       collectionAuthor.Name,
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Likes:            collectionMetrics.Likes,
			Dislikes:         collectionMetrics.Dislikes,
			Views:            collectionMetrics.Views,
			Mastered:         collectionUserProgress.Mastered,
			Reviewing:        collectionUserProgress.Reviewing,
			Learning:         collectionUserProgress.Learning,
			Starred:          collectionUserMetrics.Starred,
			IsLikedByUser:    collectionUserMetrics.Liked,
			IsDislikedByUser: collectionUserMetrics.Disliked,
			IsViewedByUser:   collectionUserMetrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collectionResponses = append(collectionResponses, collectionResponse)
	}

	return collectionResponses, err
}

func (uc *usecase) GetLikedCollectionsPreview(userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {
	collectionResponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collectionRepo.GetLikedCollectionsPreview(userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collectionAuthor, err := uc.userRepo.GetUserById(collection.AuthorId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collectionMetrics, err := uc.collectionRepo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collectionUserProgress, err := uc.collectionRepo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collectionUserMetrics, err := uc.collectionRepo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
				err = uc.collectionRepo.CreateCollectionUserMetrics(collection.Id, userId)
				if err != nil {
					logrus.Errorf("%w: %v", ErrUnexpected, err)
					return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
				}
			} else {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		}
		totalCards, err := uc.collectionRepo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collectionResponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       collectionAuthor.Name,
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collectionUserMetrics.Starred,
			Likes:            collectionMetrics.Likes,
			Dislikes:         collectionMetrics.Dislikes,
			Views:            collectionMetrics.Views,
			Mastered:         collectionUserProgress.Mastered,
			Reviewing:        collectionUserProgress.Reviewing,
			Learning:         collectionUserProgress.Learning,
			IsLikedByUser:    collectionUserMetrics.Liked,
			IsDislikedByUser: collectionUserMetrics.Disliked,
			IsViewedByUser:   collectionUserMetrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collectionResponses = append(collectionResponses, collectionResponse)
	}

	return collectionResponses, err
}

func (uc *usecase) GetStarredCollectionsPreview(userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {
	collectionResponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collectionRepo.GetStarredCollectionsPreview(userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collectionAuthor, err := uc.userRepo.GetUserById(collection.AuthorId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collectionMetrics, err := uc.collectionRepo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collectionUserProgress, err := uc.collectionRepo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collectionUserMetrics, err := uc.collectionRepo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
				err = uc.collectionRepo.CreateCollectionUserMetrics(collection.Id, userId)
				if err != nil {
					logrus.Errorf("%w: %v", ErrUnexpected, err)
					return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
				}
			} else {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		}
		totalCards, err := uc.collectionRepo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collectionResponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       collectionAuthor.Name,
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collectionUserMetrics.Starred,
			Likes:            collectionMetrics.Likes,
			Dislikes:         collectionMetrics.Dislikes,
			Views:            collectionMetrics.Views,
			Mastered:         collectionUserProgress.Mastered,
			Reviewing:        collectionUserProgress.Reviewing,
			Learning:         collectionUserProgress.Learning,
			IsLikedByUser:    collectionUserMetrics.Liked,
			IsDislikedByUser: collectionUserMetrics.Disliked,
			IsViewedByUser:   collectionUserMetrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collectionResponses = append(collectionResponses, collectionResponse)
	}

	return collectionResponses, err
}

func (uc *usecase) GetCollectionWithCards(id, userId uuid.UUID, page, size int) (*entity.GetCollectionWithCardsResponse, error) {
	var err error
	collection, err := uc.collectionRepo.GetCollection(id)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	collectionProgress, err := uc.collectionRepo.GetCollectionUserProgress(id, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	limit := size
	offset := (page - 1) * size
	cards, err := uc.collectionRepo.GetCollectionCards(id, limit, offset)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	collectionResponses := &entity.GetCollectionWithCardsResponse{
		Id:         collection.Id,
		Name:       collection.Name,
		Mastered:   collectionProgress.Mastered,
		Reviewing:  collectionProgress.Reviewing,
		Learning:   collectionProgress.Learning,
		TotalCards: cards.Total,
		Cards:      cards.Cards,
	}
	return collectionResponses, nil
}

func (uc *usecase) StarCollectionById(id, userId uuid.UUID) error {
	err := uc.collectionRepo.StarCollectionById(id, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	return nil
}

func (uc *usecase) LikeCollectionById(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error) {
	_, err := uc.collectionRepo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
			err = uc.collectionRepo.CreateCollectionUserMetrics(id, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	isLiked, isDisliked, err := uc.collectionRepo.IsCollectionLikedOrDislikedByUser(id, userId)
	if err != nil {
		// if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
		// 	return nil, ErrNotFound
		// }
		// if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
		// 	err = uc.collectionRepo.CreateCollectionUserMetrics(id, userId)
		// 	if err != nil {
		// 		logrus.Errorf("%w: %v", ErrUnexpected, err)
		// 		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		// 	}
		// }
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	if isDisliked {
		err := uc.collectionRepo.CollectionDislikeInteraction(id, userId, true)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	err = uc.collectionRepo.CollectionLikeInteraction(id, userId, isLiked)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	metrics, err := uc.collectionRepo.GetCollectionMetrics(id)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	userMetrics, err := uc.collectionRepo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
			err = uc.collectionRepo.CreateCollectionUserMetrics(id, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	return &entity.CollectionFullUserMetricsResponse{
		CollectionId: id,
		Likes:        metrics.Likes,
		Dislikes:     metrics.Dislikes,
		Views:        metrics.Views,
		UserId:       userId,
		Liked:        userMetrics.Liked,
		Disliked:     userMetrics.Disliked,
		Viewed:       userMetrics.Viewed,
		Starred:      userMetrics.Starred,
	}, nil
}

func (uc *usecase) DislikeCollectionById(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error) {
	_, err := uc.collectionRepo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			err = uc.collectionRepo.CreateCollectionUserMetrics(id, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	isLiked, isDisliked, err := uc.collectionRepo.IsCollectionLikedOrDislikedByUser(id, userId)
	if err != nil {
		// if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
		// 	return nil, ErrNotFound
		// }
		if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
			err = uc.collectionRepo.CreateCollectionUserMetrics(id, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	if isLiked {
		err := uc.collectionRepo.CollectionLikeInteraction(id, userId, true)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	err = uc.collectionRepo.CollectionDislikeInteraction(id, userId, isDisliked)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	metrics, err := uc.collectionRepo.GetCollectionMetrics(id)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	userMetrics, err := uc.collectionRepo.GetCollectionUserMetrics(id, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
			err = uc.collectionRepo.CreateCollectionUserMetrics(id, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}

	return &entity.CollectionFullUserMetricsResponse{
		CollectionId: id,
		Likes:        metrics.Likes,
		Dislikes:     metrics.Dislikes,
		Views:        metrics.Views,
		UserId:       userId,
		Liked:        userMetrics.Liked,
		Disliked:     userMetrics.Disliked,
		Viewed:       userMetrics.Viewed,
		Starred:      userMetrics.Starred,
	}, nil
}

func (uc *usecase) ViewCollectionById(id, userId uuid.UUID) error {
	isViewed, err := uc.collectionRepo.IsCollectionViewedByUser(id, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	if !isViewed {
		err := uc.collectionRepo.ViewCollection(id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}
	return nil
}

func (uc *usecase) SearchCollectionByName(text string, userId uuid.UUID) ([]*entity.UserCollectionResponse, error) {

	collectionResponses := []*entity.UserCollectionResponse{}
	var err error
	collections, err := uc.collectionRepo.SearchCollectionByName(text, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}

	for _, collection := range collections {
		collectionAuthor, err := uc.userRepo.GetUserById(collection.AuthorId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collectionMetrics, err := uc.collectionRepo.GetCollectionMetrics(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		collectionUserProgress, err := uc.collectionRepo.GetCollectionUserProgress(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}

		collectionUserMetrics, err := uc.collectionRepo.GetCollectionUserMetrics(collection.Id, userId)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionUserMetricsNotFound) {
				err = uc.collectionRepo.CreateCollectionUserMetrics(collection.Id, userId)
				if err != nil {
					logrus.Errorf("%w: %v", ErrUnexpected, err)
					return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
				}
			} else {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		}
		totalCards, err := uc.collectionRepo.GetCollectionTotal(collection.Id)
		if err != nil {
			if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
				return nil, ErrNotFound
			}
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
		createdDate := time.Date(collection.CreatedAt.Year(),
			collection.CreatedAt.Month(),
			collection.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		createdDateFormat := fmt.Sprintf("%v %v, %v", createdDate.Month(), createdDate.Day(), createdDate.Year())

		collectionResponse := &entity.UserCollectionResponse{
			Id:               collection.Id,
			Name:             collection.Name,
			AuthorName:       collectionAuthor.Name,
			Topics:           collection.Topics,
			TotalCards:       totalCards,
			Starred:          collectionUserMetrics.Starred,
			Likes:            collectionMetrics.Likes,
			Dislikes:         collectionMetrics.Dislikes,
			Views:            collectionMetrics.Views,
			Mastered:         collectionUserProgress.Mastered,
			Reviewing:        collectionUserProgress.Reviewing,
			Learning:         collectionUserProgress.Learning,
			IsLikedByUser:    collectionUserMetrics.Liked,
			IsDislikedByUser: collectionUserMetrics.Disliked,
			IsViewedByUser:   collectionUserMetrics.Viewed,
			CreatedDate:      createdDateFormat,
		}

		collectionResponses = append(collectionResponses, collectionResponse)
	}

	return collectionResponses, err

}

func (uc *usecase) CreateCollection(collection entity.Collection, cards []*entity.Card, userId uuid.UUID) error {
	createdCollection, err := uc.collectionRepo.CreateCollection(collection)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	err = uc.cardRepo.CreateMultipleCards(createdCollection.Id, cards, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
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
