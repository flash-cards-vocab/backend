package collection_usecase

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

var ErrUnexpected = errors.New("Internal error")
var ErrUnauthorized = errors.New("ErrUnauthorized")
var ErrNotFound = errors.New("ErrNotFound")
var ErrForbiddenSelfRequest = errors.New("Self request is forbidden")

type UseCase interface {
	GetMyCollections(userId uuid.UUID) ([]*entity.UserCollectionResponse, error)
	GetRecommendedCollectionsPreview(userId uuid.UUID) ([]*entity.UserCollectionResponse, error)
	GetCollectionWithCards(id, userId uuid.UUID, page, size int) (*entity.GetCollectionWithCardsResponse, error)
	StarCollectionById(id, userId uuid.UUID) error
	// GetCollectionMetrics(id, userId uuid.UUID)
	GetCollectionFullUserMetrics(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error)
	LikeCollectionById(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error)
	DislikeCollectionById(id, userId uuid.UUID) (*entity.CollectionFullUserMetricsResponse, error)
	ViewCollectionById(id, userId uuid.UUID) error
	SearchCollectionByName(text string, userId uuid.UUID) ([]*entity.UserCollectionResponse, error)
	CreateCollection(collection entity.Collection, cards []*entity.Card, userId uuid.UUID) error
	UpdateCollectionUserProgress(id uuid.UUID, mastered, reviewing, learning uint32) error
}
