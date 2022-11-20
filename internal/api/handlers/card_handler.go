package handlers

import (
	"errors"
	"net/http"

	cardUC "github.com/flash-cards-vocab/backend/app/usecase/card"
	collectionUC "github.com/flash-cards-vocab/backend/app/usecase/collection"
	handlerIntf "github.com/flash-cards-vocab/backend/internal/api/handler_interfaces"
	"github.com/flash-cards-vocab/backend/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handlerCard struct {
	cardUsecase cardUC.UseCase
	apikey      string
}

func NewCardHandler(cardUsecase cardUC.UseCase, apikey string) handlerIntf.RestCardHandler {
	return &handlerCard{cardUsecase: cardUsecase, apikey: apikey}
}

func (h *handlerCard) UploadCardImage(c *gin.Context) {
	// start_time := time.Now()

	reqApikey := c.Request.Header.Get("Apikey")
	if reqApikey != h.apikey {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "API key is incorrect"})
		return
	}

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "ERROR HERE1" + err.Error()})
		return
	}
	defer file.Close()
	storageLocation := c.Request.FormValue("location")
	if handler == nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "Invalid file"})
		return
	}
	if storageLocation == "" {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "location is required"})
		return
	}

	cloudFilename, err := h.cardUsecase.UploadCardImage(file, storageLocation, handler.Filename) // Upload(file, storageLocation, handler.Filename, random_filename != "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "here" + err.Error()})
		return
	}

	c.Request.Header.Set("Content-Type", "application/json")
	// c.JSON(http.StatusOK, resp)
	c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Data: cloudFilename})
}

func (h *handlerCard) SearchByWord(c *gin.Context) {
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}
	text := c.Param("query")

	data, err := h.cardUsecase.SearchByWord(text, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Data: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}

}

func (h *handlerCard) KnowCard(c *gin.Context) {
	paramCardId := c.Param("card_id")
	cardId, err := uuid.Parse(paramCardId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	paramCollectionId := c.Param("collection_id")
	collectionId, err := uuid.Parse(paramCollectionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}

	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	data, err := h.cardUsecase.KnowCard(collectionId, cardId, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Data: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCard) DontKnowCard(c *gin.Context) {
	paramCardId := c.Param("card_id")
	cardId, err := uuid.Parse(paramCardId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	paramCollectionId := c.Param("collection_id")
	collectionId, err := uuid.Parse(paramCollectionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}

	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	data, err := h.cardUsecase.DontKnowCard(collectionId, cardId, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Data: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCard) AddExistingCardToCollection(c *gin.Context) {
	paramCardId := c.Param("card_id")
	cardId, err := uuid.Parse(paramCardId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	paramCollectionId := c.Param("collection_id")
	collectionId, err := uuid.Parse(paramCollectionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}

	// userId, err = helpers.GetAuthContext(c)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	// }

	err = h.cardUsecase.AddExistingCardToCollection(collectionId, cardId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{"Collection Viewed"})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}
