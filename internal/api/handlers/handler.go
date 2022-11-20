package handlers

import (
	"os"

	"github.com/flash-cards-vocab/backend/app/usecase"
	handlerIntf "github.com/flash-cards-vocab/backend/internal/api/handler_interfaces"
	"github.com/flash-cards-vocab/backend/pkg/application"
)

type Handler struct {
	App               *application.Application
	CollectionHandler handlerIntf.RestCollectionHandler
	UserHandler       handlerIntf.RestUserHandler
	CardHandler       handlerIntf.RestCardHandler
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
