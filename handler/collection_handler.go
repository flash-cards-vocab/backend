package handler

import (
	"errors"
	"net/http"
	"strconv"

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

func (h *handlerCollection) GetRecommendedCollectionsPreview(c *gin.Context) {
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	data, err := h.collection_uc.GetRecommendedCollectionsPreview(user_ctx.UserId)
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

func (h *handlerCollection) GetCollectionWithCards(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil || size < 1 {
		size = 10
	}

	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	data, err := h.collection_uc.GetCollectionWithCards(id, user_ctx.UserId, page, size)
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

func (h *handlerCollection) GetCollectionMetricsById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	metrics, err := h.collection_uc.GetCollectionFullUserMetrics(id, user_ctx.UserId)

	err = h.collection_uc.StarCollectionById(id, user_ctx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{Data: metrics})
	} else {
		if errors.Is(err, collection_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}

}

func (h *handlerCollection) StarCollectionById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	err = h.collection_uc.StarCollectionById(id, user_ctx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{Data: "SUCCESS"})
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	metrics, err := h.collection_uc.LikeCollectionById(id, user_ctx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{Data: metrics})
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	metrics, err := h.collection_uc.DislikeCollectionById(id, user_ctx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{Data: metrics})
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
	if err != nil {
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
	user_ctx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	}

	collectionToCreate := entity.Collection{
		Name:     createCollectionData.Name,
		Topics:   createCollectionData.Topics,
		AuthorId: user_ctx.UserId,
	}
	err = h.collection_uc.CreateCollection(collectionToCreate, createCollectionData.Cards, user_ctx.UserId)
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
