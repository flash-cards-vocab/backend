package usecase

import (
	"context"
	"log"
	"os"

	// card_usecase "github.com/flash-cards-vocab/backend/app/usecase/card"
	"cloud.google.com/go/storage"
	card_usecase "github.com/flash-cards-vocab/backend/app/usecase/card"
	collection_usecase "github.com/flash-cards-vocab/backend/app/usecase/collection"
	user_usecase "github.com/flash-cards-vocab/backend/app/usecase/user"
	"github.com/flash-cards-vocab/backend/pkg/application"
	"github.com/flash-cards-vocab/backend/pkg/repository"
	"google.golang.org/api/option"
)

type Usecase struct {
	App               *application.Application
	UserUsecase       user_usecase.UseCase
	CollectionUsecase collection_usecase.UseCase
	CardUsecase       card_usecase.UseCase
}

func Get(app *application.Application) *Usecase {
	repo := repository.Get(app)

	api_key := `
	{
	}
	`

	gcs_client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(api_key)))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	userUsecase := user_usecase.New(repo.UserRepository)
	collectionUsecase := collection_usecase.New(repo.CollectionRepository, repo.CardRepository, repo.UserRepository)
	cardUsecase := card_usecase.New(repo.CardRepository, gcs_client, os.Getenv("GCS_BUCKET_NAME"), os.Getenv("GCS_PREFIX"))

	return &Usecase{
		App:               app,
		UserUsecase:       userUsecase,
		CollectionUsecase: collectionUsecase,
		CardUsecase:       cardUsecase,
	}
}
