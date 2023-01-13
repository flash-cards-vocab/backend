package card_usecase

import (
	"errors"
	"mime/multipart"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

var ErrUnexpected = errors.New("Internal error")
var ErrUnauthorized = errors.New("ErrUnauthorized")
var ErrNotFound = errors.New("ErrNotFound")
var ErrForbiddenSelfRequest = errors.New("Self request is forbidden")

type UseCase interface {
	// UploadCardImage(file multipart.File, location string, filename string) (string, error)
	UploadCardImage(
		file multipart.File,
		location string,
		filename string,
	) (string, error)
	AddExistingCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error
	SearchByWord(word string, userId uuid.UUID, page, size int) (*entity.CardSearch, error)
	KnowCard(collectionId, cardId, userId uuid.UUID) (*entity.CollectionUserProgress, error)
	DontKnowCard(collectionId, cardId, userId uuid.UUID) (*entity.CollectionUserProgress, error)
}
