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
	repositoryIntf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type usecase struct {
	cardRepo       repositoryIntf.CardRepository
	collectionRepo repositoryIntf.CollectionRepository
	gcsClient      *storage.Client
	bucketName     string
	envPrefix      string
}

func New(cardRepo repositoryIntf.CardRepository, collectionRepo repositoryIntf.CollectionRepository,
	gcsClient *storage.Client,
	bucketName string,
	envPrefix string,
) UseCase {
	return &usecase{
		cardRepo:       cardRepo,
		collectionRepo: collectionRepo,
		gcsClient:      gcsClient,
		bucketName:     bucketName,
		envPrefix:      envPrefix,
	}
}

func (u *usecase) UploadCardImage(
	file multipart.File,
	location string,
	filename string,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	filenameToUpload := location + "/" + strings.ReplaceAll(filename, " ", "+")
	fullFilename := "https://storage.googleapis.com/" + u.bucketName + "/" + u.envPrefix + "/" + filenameToUpload + "--" + uuid.NewString()
	wc := u.gcsClient.Bucket(u.bucketName).Object(u.envPrefix + "/" + filenameToUpload).NewWriter(ctx)
	// wc.ACL = []storage.ACLRule{{Entity: storage.AllAuthenticatedUsers, Role: storage.RoleOwner}}

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("%w: %v", "ErrUnexpected1", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("%w: %v", "ErrUnexpected2", err)
	}

	return fullFilename, nil
}

func (uc *usecase) SearchByWord(word string, userId uuid.UUID, page, size int) ([]*entity.Card, error) {
	limit := size
	offset := (page - 1) * size
	cards, err := uc.cardRepo.GetCardsByWord(word, limit, offset)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	// collectionProgress, err := uc.cardRepo.GetCollectionUserProgress(id, userId)
	// if err != nil {
	// 	if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
	// 		return nil, ErrNotFound
	// 	}
	// 	logrus.Errorf("%w: %v", ErrUnexpected, err)
	// 	return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	// }

	// limit := size
	// offset := (page - 1) * size
	// cards, err := uc.cardRepo.GetCollectionCards(id, limit, offset)
	// if err != nil {
	// 	if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
	// 		return nil, ErrNotFound
	// 	}
	// 	logrus.Errorf("%w: %v", ErrUnexpected, err)
	// 	return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	// }
	// collectionResponses := &entity.GetCollectionWithCardsResponse{
	// 	Id:         collection.Id,
	// 	Name:       collection.Name,
	// 	Mastered:   collectionProgress.Mastered,
	// 	Reviewing:  collectionProgress.Reviewing,
	// 	Learning:   collectionProgress.Learning,
	// 	TotalCards: cards.Total,
	// 	Cards:      cards.Cards,
	// }
	return cards, nil
}

func (uc *usecase) AddExistingCardToCollection(collectionId uuid.UUID, cardId uuid.UUID) error {
	err := uc.cardRepo.AssignCardToCollection(collectionId, cardId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	return nil

}

func (uc *usecase) KnowCard(collectionId, cardId, userId uuid.UUID) (*entity.CollectionUserProgress, error) {
	err := uc.cardRepo.KnowCard(collectionId, cardId, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	collUserProgr, err := uc.collectionRepo.GetCollectionUserProgress(collectionId, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionUserProgressNotFound) {
			err = uc.collectionRepo.CreateCollectionUserProgress(collectionId, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}
	return collUserProgr, nil
}

func (uc *usecase) DontKnowCard(collectionId, cardId, userId uuid.UUID) (*entity.CollectionUserProgress, error) {
	err := uc.cardRepo.DontKnowCard(collectionId, cardId, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionNotFound) {
			return nil, ErrNotFound
		}
		logrus.Errorf("%w: %v", ErrUnexpected, err)
		return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
	}
	collUserProgr, err := uc.collectionRepo.GetCollectionUserProgress(collectionId, userId)
	if err != nil {
		if errors.Is(err, repositoryIntf.ErrCollectionUserProgressNotFound) {
			err = uc.collectionRepo.CreateCollectionUserProgress(collectionId, userId)
			if err != nil {
				logrus.Errorf("%w: %v", ErrUnexpected, err)
				return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
			}
		} else {
			logrus.Errorf("%w: %v", ErrUnexpected, err)
			return nil, fmt.Errorf("%w: %v", ErrUnexpected, "Unexpected error")
		}
	}
	return collUserProgr, nil
}
