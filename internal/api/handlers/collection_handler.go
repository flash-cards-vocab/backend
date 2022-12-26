package handlers

import (
	"errors"
	"net/http"
	"strconv"

	collectionUC "github.com/flash-cards-vocab/backend/app/usecase/collection"
	"github.com/flash-cards-vocab/backend/entity"
	handlerIntf "github.com/flash-cards-vocab/backend/internal/api/handler_interfaces"
	"github.com/flash-cards-vocab/backend/pkg/helpers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handlerCollection struct {
	collectionUsecase collectionUC.UseCase
}

func NewCollectionHandler(collectionUsecase collectionUC.UseCase) handlerIntf.RestCollectionHandler {
	return &handlerCollection{collectionUsecase: collectionUsecase}
}

func (h *handlerCollection) GetMyCollections(c *gin.Context) {
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	data, err := h.collectionUsecase.GetMyCollections(userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCollection) GetCollectionUserProgress(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	data, err := h.collectionUsecase.GetCollectionUserProgress(id, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCollection) GetRecommendedCollectionsPreview(c *gin.Context) {
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil || size < 1 {
		size = 5
	}

	data, err := h.collectionUsecase.GetRecommendedCollectionsPreview(userCtx.UserId, page, size)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCollection) GetLikedCollectionsPreview(c *gin.Context) {
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	data, err := h.collectionUsecase.GetLikedCollectionsPreview(userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}

}

func (h *handlerCollection) GetStarredCollectionsPreview(c *gin.Context) {
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	data, err := h.collectionUsecase.GetStarredCollectionsPreview(userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCollection) GetCollectionWithCards(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil || size < 1 {
		size = 10
	}

	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	data, err := h.collectionUsecase.GetCollectionWithCards(id, userCtx.UserId, page, size)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCollection) GetCollectionMetricsById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	metrics, err := h.collectionUsecase.GetCollectionFullUserMetrics(id, userCtx.UserId)

	err = h.collectionUsecase.StarCollectionById(id, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: metrics})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}

}

func (h *handlerCollection) StarCollectionById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	err = h.collectionUsecase.StarCollectionById(id, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: "SUCCESS"})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCollection) LikeCollectionById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	metrics, err := h.collectionUsecase.LikeCollectionById(id, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: metrics})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}
func (h *handlerCollection) DislikeCollectionById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	metrics, err := h.collectionUsecase.DislikeCollectionById(id, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: metrics})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}
func (h *handlerCollection) ViewCollectionById(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	err = h.collectionUsecase.ViewCollectionById(id, userCtx.UserId)
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
func (h *handlerCollection) SearchCollectionByName(c *gin.Context) {
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}
	text := c.Param("query")

	data, err := h.collectionUsecase.SearchCollectionByName(text, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}

func (h *handlerCollection) CreateCollection(c *gin.Context) {
	var createCollectionData entity.CreateCollectionRequest
	err := c.ShouldBindJSON(&createCollectionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	collectionToCreate := entity.Collection{
		Name:     createCollectionData.Name,
		Topics:   createCollectionData.Topics,
		AuthorId: userCtx.UserId,
	}
	err = h.collectionUsecase.CreateCollection(collectionToCreate, createCollectionData.Cards, userCtx.UserId)
	if err == nil {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{"Collection Created"})
	} else {
		if errors.Is(err, collectionUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
}
func (h *handlerCollection) UpdateCollectionUserProgress(c *gin.Context) {
	panic("not implemented")
}

func (h *handlerCollection) UploadCollectionWithFile(c *gin.Context) {
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "ERROR HERE1" + err.Error()})
		return
	}
	defer file.Close()
	if handler == nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "Invalid file"})
		return
	}

	resp, err := h.collectionUsecase.UploadCollectionWithFile(userCtx.UserId, file, handler.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "here" + err.Error()})
		return
	}

	c.Request.Header.Set("Content-Type", "application/json")
	// c.JSON(http.StatusOK, resp)
	c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: resp})

}
