package api

import (
	"github.com/flash-cards-vocab/backend/internal/api/handlers"
	"github.com/flash-cards-vocab/backend/internal/api/middleware"
	"github.com/flash-cards-vocab/backend/pkg/application"
	"github.com/gin-gonic/gin"
)

func NewRouter(app *application.Application) (*gin.Engine, error) {
	router := gin.Default()
	h := handlers.Get(app)

	router.Use(middleware.CORSMiddleware())
	v1 := router.Group("")

	// User routes
	user := v1.Group("/user")
	user.POST("/login", h.UserHandler.Login)
	user.POST("/register", h.UserHandler.Register)
	user.GET("/profile", middleware.AuthorizeJWT, h.UserHandler.GetProfile)

	// Collection routes
	collection := v1.Group("/collection")
	collection.GET("/my", middleware.AuthorizeJWT, h.CollectionHandler.GetMyCollections)
	collection.GET("/recommended", middleware.AuthorizeJWT, h.CollectionHandler.GetRecommendedCollectionsPreview)
	collection.GET("/liked", middleware.AuthorizeJWT, h.CollectionHandler.GetLikedCollectionsPreview)
	collection.GET("/starred", middleware.AuthorizeJWT, h.CollectionHandler.GetStarredCollectionsPreview)
	collection.GET("/metrics/:id", middleware.AuthorizeJWT, h.CollectionHandler.GetCollectionMetricsById)
	collection.GET("/full/:id", middleware.AuthorizeJWT, h.CollectionHandler.GetCollectionWithCards)
	collection.PUT("/star/:id", middleware.AuthorizeJWT, h.CollectionHandler.StarCollectionById)
	collection.PUT("/like/:id", middleware.AuthorizeJWT, h.CollectionHandler.LikeCollectionById)
	collection.PUT("/dislike/:id", middleware.AuthorizeJWT, h.CollectionHandler.DislikeCollectionById)
	collection.PUT("/view/:id", middleware.AuthorizeJWT, h.CollectionHandler.ViewCollectionById)
	collection.GET("/search/:query", middleware.AuthorizeJWT, h.CollectionHandler.SearchCollectionByName)
	collection.POST("/create", middleware.AuthorizeJWT, h.CollectionHandler.CreateCollection)
	collection.PUT("/update-user-progress/:id", middleware.AuthorizeJWT, h.CollectionHandler.UpdateCollectionUserProgress)

	// Card routes
	card := v1.Group("/card")
	card.POST("/upload-card-image", middleware.AuthorizeJWT, h.CardHandler.UploadCardImage)
	card.POST("/add-card-to-collection/:collection_id/:card_id", middleware.AuthorizeJWT, h.CardHandler.AddExistingCardToCollection)
	card.PUT("/know/:card_id/:collection_id", middleware.AuthorizeJWT, h.CardHandler.KnowCard)
	card.PUT("/dont-know/:card_id/:collection_id", middleware.AuthorizeJWT, h.CardHandler.DontKnowCard)

	return router, nil
}
