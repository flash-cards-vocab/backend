package repository

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
)

var ErrUserNotFound = errors.New("loan request not found")

type UserRepository interface {
	// GetMyCollectionsPreview(user_id uuid.UUID) (entity.Collection, error)
	// LikeCollection(employee_id string, limit, offset int) ([]entity.Collection, error)
	// DislikeCollection(id int) (*entity.Collection, error)
	// ViewCollection(id int) (*entity.Collection, error)
	// SearchCollectionByName(id int, attachment string) error
	// CreateCollection(collectionName string, collectionCards entity.Card) error

	CreateUser(user entity.User) (*entity.User, error)
	CheckIfUserExistsByEmail(email string) (bool, error)
	GetUserByEmail(email string) (*entity.User, error)
}
