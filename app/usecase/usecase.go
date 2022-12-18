package usecase

import (
	"context"
	"log"
	"os"

	// card_usecase "github.com/flash-cards-vocab/backend/app/usecase/card"
	"cloud.google.com/go/storage"
	cardUC "github.com/flash-cards-vocab/backend/app/usecase/card"
	collectionUC "github.com/flash-cards-vocab/backend/app/usecase/collection"
	userUC "github.com/flash-cards-vocab/backend/app/usecase/user"
	"github.com/flash-cards-vocab/backend/pkg/application"
	"github.com/flash-cards-vocab/backend/pkg/repository"
	"google.golang.org/api/option"
)

type Usecase struct {
	App               *application.Application
	UserUsecase       userUC.UseCase
	CollectionUsecase collectionUC.UseCase
	CardUsecase       cardUC.UseCase
}

func Get(app *application.Application) *Usecase {
	repo := repository.Get(app)

	apiKey := `
	{
	}
	`

	gcsClient, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(apiKey)))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	userUsecase := userUC.New(repo.UserRepository, repo.CompanyRepository)
	collectionUsecase := collectionUC.New(repo.CollectionRepository, repo.CardRepository, repo.UserRepository)
	cardUsecase := cardUC.New(repo.CardRepository, repo.CollectionRepository, gcsClient, os.Getenv("GCS_BUCKET_NAME"), os.Getenv("GCS_PREFIX"))

	return &Usecase{
		App:               app,
		UserUsecase:       userUsecase,
		CollectionUsecase: collectionUsecase,
		CardUsecase:       cardUsecase,
	}
}
