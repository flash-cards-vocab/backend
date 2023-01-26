package collection_usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	repositoryIntf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type usecase struct {
	collectionRepo repositoryIntf.CollectionRepository
	cardRepo       repositoryIntf.CardRepository
	userRepo       repositoryIntf.UserRepository
	gcsClient      *storage.Client
	bucketName     string
	envPrefix      string
}

func New(
	collectionRepo repositoryIntf.CollectionRepository,
	cardRepo repositoryIntf.CardRepository,
	userRepo repositoryIntf.UserRepository,
	gcsClient *storage.Client,
	bucketName string,
	envPrefix string,
) UseCase {
	return &usecase{
		collectionRepo: collectionRepo,
		cardRepo:       cardRepo,
		userRepo:       userRepo,
		gcsClient:      gcsClient,
		bucketName:     bucketName,
		envPrefix:      envPrefix,
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

func (uc *usecase) GetRecommendedCollectionsPreview(userId uuid.UUID, page, size int) ([]*entity.UserCollectionResponse, error) {
	collectionResponses := []*entity.UserCollectionResponse{}
	// var err error
	limit := size
	offset := (page - 1) * size

	collections, err := uc.collectionRepo.GetRecommendedCollectionsPreview(userId, limit, offset)
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

func (uc *usecase) GetCollectionWithCards(collectionId, userId uuid.UUID, page, size int) (*entity.GetCollectionWithCardsResponse, error) {
	var err error
	collection, err := uc.collectionRepo.GetCollection(collectionId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	collectionProgress, err := uc.collectionRepo.GetCollectionUserProgress(collectionId, userId)
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
	limit := size
	offset := (page - 1) * size
	cards, err := uc.collectionRepo.GetCollectionCards(collectionId, userId, limit, offset)
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
		Topics:     collection.Topics,
		Cards:      cards.CardForUser,
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
	urlGCP := "https://storage.googleapis.com/flashcards-images"

	for _, card := range cards {
		if !strings.Contains(card.ImageUrl, urlGCP) {
			response, err := http.Get(card.ImageUrl)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			if response.StatusCode != 200 {
				return fmt.Errorf("%w: %v", ErrUnexpected, "Received non 200 response code")
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
			defer cancel()

			filenameToUpload := "card_images/" + strings.ReplaceAll(card.Word, " ", "+") + "--" + uuid.NewString()
			fileURL := "https://storage.googleapis.com/" + uc.bucketName + "/" + uc.envPrefix + "/" + filenameToUpload
			wc := uc.gcsClient.Bucket(uc.bucketName).Object(uc.envPrefix + "/" + filenameToUpload).NewWriter(ctx)
			// wc.ACL = []storage.ACLRule{{Entity: storage.AllAuthenticatedUsers, Role: storage.RoleOwner}}

			if _, err := io.Copy(wc, response.Body); err != nil {
				return fmt.Errorf("%w: %v", "ErrUnexpected1", err)
			}

			if err := wc.Close(); err != nil {
				return fmt.Errorf("%w: %v", "ErrUnexpected2", err)
			}
			card.ImageUrl = fileURL

		}
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

func (uc *usecase) UploadCollectionWithFile(userId uuid.UUID, file multipart.File, filename string) (*entity.CreateMultipleCollectionResponse, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer func() error {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
		return err
	}()

	collectionEnt := entity.Collection{}
	cards := []*entity.Card{}

	collectionName, err := f.GetCellValue(f.GetSheetName(0), "B3")
	if err != nil {
		return nil, err
	}
	collectionTopics, err := f.GetCellValue(f.GetSheetName(0), "B4")
	if err != nil {
		return nil, err
	}
	collectionEnt.Name = collectionName
	collectionEnt.Topics = strings.Split(collectionTopics, ";")
	collectionEnt.AuthorId = userId

	collection, err := uc.collectionRepo.CreateCollection(collectionEnt)

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rowI := 7; rowI < len(rows); rowI++ {
		fmt.Println(rowI)
		// Skip row that does not contain one of mandatory fields
		if len(rows[rowI]) < 4 {
			continue
		}
		card := &entity.Card{
			Id:         uuid.New(),
			Word:       rows[rowI][0],
			Definition: rows[rowI][1],
			Sentence:   rows[rowI][2],
			ImageUrl:   rows[rowI][3],
		}

		if len(rows[rowI]) >= 5 {
			card.Antonyms = rows[rowI][4]
		}
		if len(rows[rowI]) >= 6 {
			card.Synonyms = rows[rowI][5]
		}
		card.AuthorId = userId

		cards = append(cards, card)
	}
	err = uc.cardRepo.CreateMultipleCards(collection.Id, cards, userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &entity.CreateMultipleCollectionResponse{
		Name:        collection.Name,
		CardsAmount: uint32(len(cards)),
	}, nil
}

func (uc *usecase) UpdateCollection(
	userId uuid.UUID,
	updateData *entity.UpdateCollectionRequest) error {

	collectionData := entity.Collection{
		Id:     updateData.Id,
		Name:   updateData.Name,
		Topics: updateData.Topics,
	}
	err := uc.collectionRepo.UpdateCollection(collectionData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cardsToCreate := []*entity.Card{}
	cardsToRemove := []*entity.CollectionCards{}

	for _, card := range updateData.Cards {
		switch card.Action {
		case entity.CardUpdateType_Create: // just create a new card
			card := &entity.Card{
				Word:       card.Word,
				ImageUrl:   card.ImageUrl,
				Definition: card.Definition,
				Sentence:   card.Sentence,
				Antonyms:   card.Antonyms,
				Synonyms:   card.Synonyms,
			}
			cardsToCreate = append(cardsToCreate, card)

		case entity.CardUpdateType_Update: // create a new card and remove the previous from the collection
			cardToCreate := &entity.Card{
				Word:       card.Word,
				ImageUrl:   card.ImageUrl,
				Definition: card.Definition,
				Sentence:   card.Sentence,
				Antonyms:   card.Antonyms,
				Synonyms:   card.Synonyms,
			}
			cardsToCreate = append(cardsToCreate, cardToCreate)

			cardToRemove := &entity.CollectionCards{
				CardId:       card.Id,
				CollectionId: collectionData.Id,
			}
			cardsToRemove = append(cardsToRemove, cardToRemove)

		case entity.CardUpdateType_Remove: // remove the card from the collection, do not delete it
			card := &entity.CollectionCards{
				CardId:       card.Id,
				CollectionId: collectionData.Id,
			}
			cardsToRemove = append(cardsToRemove, card)

		}
	}
	if len(cardsToCreate) > 0 {
		err = uc.cardRepo.CreateMultipleCards(collectionData.Id, cardsToCreate, userId)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	if len(cardsToRemove) > 0 {
		err = uc.cardRepo.RemoveMultipleCardsFromCollection(cardsToRemove)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
