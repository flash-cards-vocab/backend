package handlers

import (
	"os"

	"github.com/flash-cards-vocab/backend/app/usecase"
	"github.com/flash-cards-vocab/backend/internal/api/handler_interfaces"
	"github.com/flash-cards-vocab/backend/pkg/application"
)

type Handler struct {
	App               *application.Application
	CollectionHandler handler_interfaces.RestCollectionHandler
	UserHandler       handler_interfaces.RestUserHandler
	CardHandler       handler_interfaces.RestCardHandler
}

func Get(app *application.Application) *Handler {
	uc := usecase.Get(app)
	userHandler := NewUserHandler(uc.UserUsecase)
	collectionHandler := NewCollectionHandler(uc.CollectionUsecase)
	cardHandler := NewCardHandler(uc.CardUsecase, os.Getenv("GCS_API_KEY"))

	return &Handler{
		App:               app,
		UserHandler:       userHandler,
		CollectionHandler: collectionHandler,
		CardHandler:       cardHandler,
	}
}
