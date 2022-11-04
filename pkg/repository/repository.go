package repository

import (
	repositoryIntf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/pkg/application"
	"github.com/flash-cards-vocab/backend/pkg/repository/card_repository"
	"github.com/flash-cards-vocab/backend/pkg/repository/collection_repository"
	"github.com/flash-cards-vocab/backend/pkg/repository/user_repository"
)

type Repository struct {
	CardRepository       repositoryIntf.CardRepository
	CollectionRepository repositoryIntf.CollectionRepository
	UserRepository       repositoryIntf.UserRepository
}

func Get(app *application.Application) *Repository {
	cardRepository := card_repository.New(app.DBManager.DB)
	collectionRepository := collection_repository.New(app.DBManager.DB)
	userRepository := user_repository.New(app.DBManager.DB)

	return &Repository{
		CardRepository:       cardRepository,
		CollectionRepository: collectionRepository,
		UserRepository:       userRepository,
	}
}
