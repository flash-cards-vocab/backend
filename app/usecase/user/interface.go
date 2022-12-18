package user_usecase

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
)

var ErrUnexpected = errors.New("Internal error")
var ErrUserExistsAlready = errors.New("User exists already")
var ErrUnauthorized = errors.New("ErrUnauthorized")
var ErrNotFound = errors.New("ErrNotFound")
var ErrForbiddenSelfRequest = errors.New("Self request is forbidden")
var ErrUserPasswordMismatch = errors.New("User password is Incorrect")

type UseCase interface {
	Register(user entity.UserRegistration) (*entity.UserWithAuthToken, error)
	Login(user entity.UserLogin) (*entity.UserWithAuthToken, error)
}
