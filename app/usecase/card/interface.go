package card_usecase

import (
	"errors"
	"mime/multipart"

	"github.com/google/uuid"
)

var ErrUnexpected = errors.New("Internal error")
var ErrUnauthorized = errors.New("Anda tidak memiliki akses")
var ErrNotFound = errors.New("Permintaan pinjaman tidak ditemukan")
var ErrForbiddenSelfRequest = errors.New("Self request is forbidden")

type UseCase interface {
	// UploadCardImage(file multipart.File, location string, filename string) (string, error)
	UploadCardImage(
		file multipart.File,
		location string,
		filename string,
	) (string, error)
	AddExistingCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error
}
