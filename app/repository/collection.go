package repository

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

var ErrCollectionNotFound = errors.New("loan request not found")
var ErrCollectionUserMetricsNotFound = errors.New("collection user metrics not found")

type CollectionRepository interface {
	GetMyCollections(userId uuid.UUID) ([]*entity.Collection, error)
	GetCollectionTotal(collection_id uuid.UUID) (int, error)
	GetRecommendedCollectionsPreview(userId uuid.UUID) ([]*entity.Collection, error)
	GetLikedCollectionsPreview(userId uuid.UUID) ([]*entity.Collection, error)
	GetStarredCollectionsPreview(userId uuid.UUID) ([]*entity.Collection, error)
	IsCollectionLikedOrDislikedByUser(id, userId uuid.UUID) (bool, bool, error)

	IsCollectionLikedByUser(id, userId uuid.UUID) (bool, error)
	IsCollectionDislikedByUser(id, userId uuid.UUID) (bool, error)
	IsCollectionViewedByUser(id, userId uuid.UUID) (bool, error)
	StarCollectionById(id, userId uuid.UUID) error
	CollectionLikeInteraction(id, userId uuid.UUID, isLiked bool) error
	CollectionDislikeInteraction(id, userId uuid.UUID, isDisliked bool) error
	ViewCollection(id, userId uuid.UUID) error
	SearchCollectionByName(search string, userId uuid.UUID) ([]*entity.Collection, error)
	CreateCollection(collection entity.Collection) (*entity.Collection, error)

	GetCollectionMetrics(id uuid.UUID) (*entity.CollectionMetrics, error)
	GetCollectionUserProgress(id, userId uuid.UUID) (*entity.CollectionUserProgress, error)
	GetCollectionUserMetrics(id, userId uuid.UUID) (*entity.CollectionUserMetrics, error)
	CreateCollectionUserMetrics(id, userId uuid.UUID) error
	CreateCollectionUserProgress(id, userId uuid.UUID) error
	GetCollection(id uuid.UUID) (*entity.Collection, error)
	GetCollectionCards(id uuid.UUID, limit, offset int) (*entity.CardsPagination, error)

	// CreateCollectionUserInteraction()
}
