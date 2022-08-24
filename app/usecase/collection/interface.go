package collection_usecase

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

var ErrUnexpected = errors.New("Internal error")
var ErrUnauthorized = errors.New("Anda tidak memiliki akses")
var ErrNotFound = errors.New("Permintaan pinjaman tidak ditemukan")
var ErrForbiddenSelfRequest = errors.New("Self request is forbidden")

type UseCase interface {
	GetMyCollections(userId uuid.UUID) ([]*entity.UserCollectionResponse, error)
	LikeCollectionById(id, userId uuid.UUID) error
	DislikeCollectionById(id, userId uuid.UUID) error
	ViewCollectionById(id, userId uuid.UUID) error
	SearchCollectionByName(text string) ([]*entity.Collection, error)
	CreateCollection(collection entity.Collection, cards []*entity.Card) error
	UpdateCollectionUserProgress(id uuid.UUID, mastered, reviewing, learning uint32) error
}
