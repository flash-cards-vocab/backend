package repository

import (
	repositoryIntf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/pkg/application"
	cardRepo "github.com/flash-cards-vocab/backend/pkg/repository/card_repository"
	collectionRepo "github.com/flash-cards-vocab/backend/pkg/repository/collection_repository"
	companyRepo "github.com/flash-cards-vocab/backend/pkg/repository/company_repository"
	userRepo "github.com/flash-cards-vocab/backend/pkg/repository/user_repository"
)

type Repository struct {
	CardRepository       repositoryIntf.CardRepository
	CollectionRepository repositoryIntf.CollectionRepository
	UserRepository       repositoryIntf.UserRepository
	CompanyRepository    repositoryIntf.CompanyRepository
}

func Get(app *application.Application) *Repository {
	cardRepository := cardRepo.New(app.DBManager.DB)
	collectionRepository := collectionRepo.New(app.DBManager.DB)
	userRepository := userRepo.New(app.DBManager.DB)
	companyRepository := companyRepo.New(app.DBManager.DB)

	return &Repository{
		CardRepository:       cardRepository,
		CollectionRepository: collectionRepository,
		UserRepository:       userRepository,
		CompanyRepository:    companyRepository,
	}
}
