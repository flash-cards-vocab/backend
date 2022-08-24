package handler

import (
	"errors"
	"net/http"

	collection_usecase "github.com/flash-cards-vocab/backend/app/usecase/collection"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/flash-cards-vocab/backend/pkg/helpers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handlerCollection struct {
	collection_uc collection_usecase.UseCase
}

func NewCollectionHandler(collection_uc collection_usecase.UseCase) RestCollectionHandler {
	return &handlerCollection{collection_uc: collection_uc}
}

func (h *handlerCollection) GetMyCollections(c *gin.Context) {
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	data, err := h.collection_uc.GetMyCollections(user_ctx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{Data: data})
	} else {
		if errors.Is(err, collection_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}
}
func (h *handlerCollection) LikeCollectionById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	err = h.collection_uc.LikeCollectionById(id, user_ctx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{"Collection Liked"})
	} else {
		if errors.Is(err, collection_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}
}
func (h *handlerCollection) DislikeCollectionById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	err = h.collection_uc.DislikeCollectionById(id, user_ctx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{"Collection Disliked"})
	} else {
		if errors.Is(err, collection_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}
}
func (h *handlerCollection) ViewCollectionById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	err = h.collection_uc.ViewCollectionById(id, user_ctx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{"Collection Viewed"})
	} else {
		if errors.Is(err, collection_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}
}
func (h *handlerCollection) SearchCollectionByName(c *gin.Context) {
	panic("not implemented")
}
func (h *handlerCollection) CreateCollection(c *gin.Context) {
	var createCollectionData entity.CreateCollectionRequest
	err := c.ShouldBindJSON(&createCollectionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	err = h.collection_uc.CreateCollection(entity.Collection{Name: createCollectionData.Name, Topics: createCollectionData.Topics}, createCollectionData.Cards)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{"Collection Created"})
	} else {
		if errors.Is(err, collection_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}
}
func (h *handlerCollection) UpdateCollectionUserProgress(c *gin.Context) {
	panic("not implemented")
}
