package api

import (
	"github.com/flash-cards-vocab/backend/internal/api/handlers"
	"github.com/flash-cards-vocab/backend/internal/api/middleware"
	"github.com/flash-cards-vocab/backend/pkg/application"
	"github.com/gin-gonic/gin"
)

func NewRouter(app *application.Application) (*gin.Engine, error) {
	// gin.SetMode((gin.ReleaseMode))
	router := gin.Default()
	h := handlers.Get(app)

	router.Use(middleware.CORSMiddleware())
	v1 := router.Group("")

	// User routes
	user := v1.Group("/user")
	// User GET requests
	user.GET("/profile", middleware.AuthorizeJWT, h.UserHandler.GetProfile)
	// User POST requests
	user.POST("/login", h.UserHandler.Login)
	user.POST("/register", h.UserHandler.Register)
	user.GET("/register/username-exists/:username", h.UserHandler.UsernameExists)

	// Collection routes
	collection := v1.Group("/collection")
	// Collection GET requests
	collection.GET("/my", middleware.AuthorizeJWT, h.CollectionHandler.GetMyCollections)
	collection.GET("/recommended", middleware.AuthorizeJWT, h.CollectionHandler.GetRecommendedCollectionsPreview)
	collection.GET("/liked", middleware.AuthorizeJWT, h.CollectionHandler.GetLikedCollectionsPreview)
	collection.GET("/starred", middleware.AuthorizeJWT, h.CollectionHandler.GetStarredCollectionsPreview)
	collection.GET("/metrics/:id", middleware.AuthorizeJWT, h.CollectionHandler.GetCollectionMetricsById)
	collection.GET("/full/:id", middleware.AuthorizeJWT, h.CollectionHandler.GetCollectionWithCards)
	collection.GET("/search/:query", middleware.AuthorizeJWT, h.CollectionHandler.SearchCollectionByName)
	// Collection POST requests
	collection.POST("/create", middleware.AuthorizeJWT, h.CollectionHandler.CreateCollection)
	collection.POST("/upload-collection-with-file", middleware.AuthorizeJWT, h.CollectionHandler.UploadCollectionWithFile)
	// Collection PUT requests
	collection.PUT("/update-user-progress/:id", middleware.AuthorizeJWT, h.CollectionHandler.UpdateCollectionUserProgress)
	collection.PUT("/star/:id", middleware.AuthorizeJWT, h.CollectionHandler.StarCollectionById)
	collection.PUT("/like/:id", middleware.AuthorizeJWT, h.CollectionHandler.LikeCollectionById)
	collection.PUT("/dislike/:id", middleware.AuthorizeJWT, h.CollectionHandler.DislikeCollectionById)
	collection.PUT("/view/:id", middleware.AuthorizeJWT, h.CollectionHandler.ViewCollectionById)
	collection.PUT("/update", middleware.AuthorizeJWT, h.CollectionHandler.UpdateCollection)

	// Card routes
	card := v1.Group("/card")
	// Card GET requests
	card.GET("/search-by-word", middleware.AuthorizeJWT, h.CardHandler.SearchByWord)
	// Card POST requests
	card.POST("/upload-card-image", middleware.AuthorizeJWT, h.CardHandler.UploadCardImage)
	card.POST("/add-card-to-collection/:collection_id/:card_id", middleware.AuthorizeJWT, h.CardHandler.AddExistingCardToCollection)
	// Card PUT requests
	card.PUT("/know/:card_id/:collection_id", middleware.AuthorizeJWT, h.CardHandler.KnowCard)
	card.PUT("/dont-know/:card_id/:collection_id", middleware.AuthorizeJWT, h.CardHandler.DontKnowCard)

	// Open routes
	unregistered := v1.Group("/unregistered")
	// Open collection routes
	collectionUnregistered := unregistered.Group("/collection")
	collectionUnregistered.GET("/recommended", h.CollectionHandler.UnregisteredGetRecommendedCollectionsPreview)
	collectionUnregistered.GET("/full/:id", h.CollectionHandler.UnregisteredGetCollectionWithCards)
	collectionUnregistered.GET("/search/:query", h.CollectionHandler.UnregisteredSearchCollectionByName)

	// Open card routes
	// cardUnregistered := unregistered.Group("/card")

	return router, nil
}
