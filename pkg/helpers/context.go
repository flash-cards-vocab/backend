package helpers

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAuthContext(c *gin.Context) (*entity.AuthContext, error) {
	userCtx, ok := c.Get("Id")
	if !ok {
		return nil, errors.New("user Id not found in context")
	}
	userId, err := uuid.Parse(userCtx.(string))
	if err != nil {
		return nil, errors.New("user Id is not a uuid format")
	}
	emailCtx, ok := c.Get("Email")
	if !ok {
		return nil, errors.New("email not found in context")
	}

	return &entity.AuthContext{
		UserId: userId,
		Email:  emailCtx.(string),
	}, nil
}
