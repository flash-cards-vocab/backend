package card_usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/flash-cards-vocab/backend/app/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type usecase struct {
	card_repo   repository.CardRepository
	gcs_client  *storage.Client
	bucket_name string
	env_prefix  string
}

func New(card_repo repository.CardRepository,
	gcs_client *storage.Client,
	bucket_name string,
	env_prefix string,
) UseCase {
	return &usecase{
		card_repo:   card_repo,
		gcs_client:  gcs_client,
		bucket_name: bucket_name,
		env_prefix:  env_prefix,
	}
}

func (u *usecase) UploadCardImage(
	file multipart.File,
	location string,
	filename string,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	filename_to_upload := location + "/" + strings.ReplaceAll(filename, " ", "+")
	full_filename := "https://storage.googleapis.com/" + u.bucket_name + "/" + u.env_prefix + "/" + filename_to_upload
	fmt.Println(location, filename, "location , filename")
	fmt.Println(full_filename, "full_filename")
	wc := u.gcs_client.Bucket(u.bucket_name).Object(u.env_prefix + "/" + filename_to_upload).NewWriter(ctx)
	// wc.ACL = []storage.ACLRule{{Entity: storage.AllAuthenticatedUsers, Role: storage.RoleOwner}}

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("%w: %v", "ErrUnexpected1", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("%w: %v", "ErrUnexpected2", err)
	}

	return full_filename, nil
}

func (uc *usecase) AddExistingCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error {
	err := uc.card_repo.AssignCardToCollection(collectionId, cardId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	return nil

}

func (uc *usecase) KnowCard(collectionId, cardId, userId uuid.UUID) error {
	err := uc.card_repo.KnowCard(collectionId, cardId, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	return nil
}

func (uc *usecase) DontKnowCard(collectionId, cardId, userId uuid.UUID) error {
	err := uc.card_repo.DontKnowCard(collectionId, cardId, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	return nil
}
