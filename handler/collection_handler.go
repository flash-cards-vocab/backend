package handler

import (
	"errors"
	"net/http"

	collection_usecase "github.com/flash-cards-vocab/backend/app/usecase/collection"
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
	paramId := c.Param("user_id")
	user_id, err := uuid.Parse(paramId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	data, err := h.collection_uc.GetMyCollections(user_id)
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
	paramId = c.Param("user_id")
	user_id, err := uuid.Parse(paramId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	err = h.collection_uc.LikeCollectionById(id, user_id)
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
	paramId = c.Param("user_id")
	user_id, err := uuid.Parse(paramId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	err = h.collection_uc.DislikeCollectionById(id, user_id)
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
	paramId = c.Param("user_id")
	user_id, err := uuid.Parse(paramId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	err = h.collection_uc.ViewCollectionById(id, user_id)
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
	panic("not implemented")
}
func (h *handlerCollection) UpdateCollectionUserProgress(c *gin.Context) {
	panic("not implemented")
}
