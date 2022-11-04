package repository

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("loan request not found")

type UserRepository interface {
	CreateUser(user entity.User) (*entity.User, error)
	CheckIfUserExistsByEmail(email string) (bool, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserById(id uuid.UUID) (*entity.User, error)
}
