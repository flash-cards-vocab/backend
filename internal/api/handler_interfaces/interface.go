package handler_interfaces

import (
	"github.com/gin-gonic/gin"
)

type RestCollectionHandler interface {
	GetMyCollections(c *gin.Context)
	GetRecommendedCollectionsPreview(c *gin.Context)
	GetLikedCollectionsPreview(c *gin.Context)
	GetStarredCollectionsPreview(c *gin.Context)
	GetCollectionMetricsById(c *gin.Context)
	GetCollectionWithCards(c *gin.Context)
	LikeCollectionById(c *gin.Context)
	DislikeCollectionById(c *gin.Context)
	ViewCollectionById(c *gin.Context)
	SearchCollectionByName(c *gin.Context)
	CreateCollection(c *gin.Context)
	UpdateCollectionUserProgress(c *gin.Context)
	StarCollectionById(c *gin.Context)
	GetCollectionUserProgress(c *gin.Context)
}

type RestUserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type RestCardHandler interface {
	AddExistingCardToCollection(c *gin.Context)
	UploadCardImage(c *gin.Context)
	KnowCard(c *gin.Context)
	DontKnowCard(c *gin.Context)
}
