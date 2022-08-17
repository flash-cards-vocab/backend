package user_usecase

import (
	"context"
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
)

var ErrUnexpected = errors.New("Internal error")
var ErrUserExistsAlready = errors.New("User exists already")
var ErrUnauthorized = errors.New("Anda tidak memiliki akses")
var ErrNotFound = errors.New("Permintaan pinjaman tidak ditemukan")
var ErrForbiddenSelfRequest = errors.New("Self request is forbidden")

type UseCase interface {
	Register(ctx context.Context, user entity.User) (*entity.UserWithToken, error)
	Login(ctx context.Context, user entity.User) (*entity.UserWithToken, error)
}
