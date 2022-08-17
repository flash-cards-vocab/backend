package handler

import (
	"github.com/gin-gonic/gin"
)

type RestCollectionHandler interface {
	GetMyCollections(c *gin.Context)
	LikeCollectionById(c *gin.Context)
	DislikeCollectionById(c *gin.Context)
	ViewCollectionById(c *gin.Context)
	SearchCollectionByName(c *gin.Context)
	CreateCollection(c *gin.Context)
	UpdateCollectionUserProgress(c *gin.Context)
}

type RestUserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}
