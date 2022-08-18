package handler

import (
	"errors"
	"net/http"

	user_usecase "github.com/flash-cards-vocab/backend/app/usecase/user"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/gin-gonic/gin"
)

type handlerUser struct {
	user_uc user_usecase.UseCase
}

func NewUserHandler(user_uc user_usecase.UseCase) RestUserHandler {
	return &handlerUser{user_uc: user_uc}
}

func (h *handlerUser) Register(c *gin.Context) {
	// paramId := c.Param("user_id")
	// user_id, err := uuid.Parse(paramId)
	// if err == nil {
	// 	c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	// }
	var newUserData entity.User
	err := c.ShouldBindJSON(&newUserData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	data, err := h.user_uc.Register(c, newUserData)
	if err != nil {
		if errors.Is(err, user_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

func (h *handlerUser) Login(c *gin.Context) {
	// paramId := c.Param("user_id")
	// user_id, err := uuid.Parse(paramId)
	// if err == nil {
	// 	c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	// }
	var loginUserData entity.UserLogin
	err := c.ShouldBindJSON(&loginUserData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	data, err := h.user_uc.Login(c, loginUserData)
	if err != nil {
		if errors.Is(err, user_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}
