package repository

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

var ErrCollectionNotFound = errors.New("loan request not found")

type CollectionRepository interface {
	GetMyCollections(userId uuid.UUID) ([]*entity.Collection, error)
	IsCollectionLikedByUser(id, userId uuid.UUID) (bool, error)
	IsCollectionDislikedByUser(id, userId uuid.UUID) (bool, error)
	IsCollectionViewedByUser(id, userId uuid.UUID) (bool, error)
	CollectionLikeInteraction(id, userId uuid.UUID, isLiked bool) error
	CollectionDislikeInteraction(id, userId uuid.UUID, isDisliked bool) error
	ViewCollection(id, userId uuid.UUID) error
	SearchCollectionByName(name string) error
	CreateCollection(collection entity.Collection) (*entity.Collection, error)
}
