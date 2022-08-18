package helpers

import (
	"errors"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAuthContext(c *gin.Context) (*entity.AuthContext, error) {
	user_ctx, ok := c.Get("Id")
	if !ok {
		return nil, errors.New("user Id not found in context")
	}
	user_id, err := uuid.Parse(user_ctx.(string))
	if err != nil {
		return nil, errors.New("user Id is not a uuid format")
	}
	email_ctx, ok := c.Get("Email")
	if !ok {
		return nil, errors.New("email not found in context")
	}

	return &entity.AuthContext{
		UserId: user_id,
		Email:  email_ctx.(string),
	}, nil
}
