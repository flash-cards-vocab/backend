package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	card_uc "github.com/flash-cards-vocab/backend/app/usecase/card"
	collection_uc "github.com/flash-cards-vocab/backend/app/usecase/collection"
	user_uc "github.com/flash-cards-vocab/backend/app/usecase/user"
	"github.com/flash-cards-vocab/backend/config"
	"github.com/flash-cards-vocab/backend/handler"
	"github.com/flash-cards-vocab/backend/middleware"
	"github.com/flash-cards-vocab/backend/repository/card_repository"
	"github.com/flash-cards-vocab/backend/repository/collection_repository"
	"github.com/flash-cards-vocab/backend/repository/user_repository"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"gorm.io/gorm/logger"

	postgres "github.com/flash-cards-vocab/backend/pkg/database"
)

func main() {
	cfg := config.New()
	router := gin.Default()

	db, err := postgres.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUserName, cfg.DBPassword, cfg.DBDatabaseName, logger.LogLevel(cfg.DBLogMode))
	if err != nil {
		log.Panicln("Failed to Initialized postgres DB:", err)
	}
	// rdsPool, err := redis.Connect(cfg.RedisHost, cfg.RedisPort, "")
	// if err != nil {
	// 	log.Panicln("Failed to Initialized redis:", err)
	// }

	router.Use(middleware.CORSMiddleware())

	user_repo := user_repository.New(db)
	user_uc := user_uc.New(user_repo)

	userHndlr := handler.NewUserHandler(user_uc)
	user := router.Group("/user")
	user.POST("/login", userHndlr.Login)
	user.POST("/register", userHndlr.Register)
	// user.PUT("/logout", userHndlr.DislikeCollectionById)
	// user.PUT("/refresh-token", userHndlr.ViewCollectionById)

	collection_repo := collection_repository.New(db)
	card_repo := card_repository.New(db)
	collection_uc := collection_uc.New(collection_repo, card_repo, user_repo)

	collectionHndlr := handler.NewCollectionHandler(collection_uc)

	/* CRUD author for CMS */
	collection := router.Group("/collection")
	collection.GET("/my", middleware.AuthorizeJWT, collectionHndlr.GetMyCollections)
	collection.GET("/recommended", middleware.AuthorizeJWT, collectionHndlr.GetRecommendedCollectionsPreview)
	collection.GET("/liked", middleware.AuthorizeJWT, collectionHndlr.GetLikedCollectionsPreview)
	collection.GET("/starred", middleware.AuthorizeJWT, collectionHndlr.GetStarredCollectionsPreview)
	collection.GET("/metrics/:id", middleware.AuthorizeJWT, collectionHndlr.GetCollectionMetricsById)
	collection.GET("/full/:id", middleware.AuthorizeJWT, collectionHndlr.GetCollectionWithCards)
	collection.PUT("/star/:id", middleware.AuthorizeJWT, collectionHndlr.StarCollectionById)
	collection.PUT("/like/:id", middleware.AuthorizeJWT, collectionHndlr.LikeCollectionById)
	collection.PUT("/dislike/:id", middleware.AuthorizeJWT, collectionHndlr.DislikeCollectionById)
	collection.PUT("/view/:id", middleware.AuthorizeJWT, collectionHndlr.ViewCollectionById)
	collection.GET("/search/:query", middleware.AuthorizeJWT, collectionHndlr.SearchCollectionByName)
	collection.POST("/create", middleware.AuthorizeJWT, collectionHndlr.CreateCollection)
	collection.PUT("/update-user-progress/:id", middleware.AuthorizeJWT, collectionHndlr.UpdateCollectionUserProgress)

	api_key := cfg.GCSJSONAPIKey
	gcs_client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(api_key)))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	card_uc := card_uc.New(card_repo, gcs_client, cfg.GCSBucketName, cfg.GCSPrefix)
	cardHndlr := handler.NewCardHandler(card_uc, cfg.GCSAPIKey)
	card := router.Group("/card")
	card.POST("/upload-card-image", middleware.AuthorizeJWT, cardHndlr.UploadCardImage)
	card.POST("/add-card-to-collection/:collection_id/:card_id", middleware.AuthorizeJWT, cardHndlr.AddExistingCardToCollection)
	card.PUT("/know/:card_id/:collection_id", middleware.AuthorizeJWT, cardHndlr.KnowCard)
	card.PUT("/dont-know/:card_id/:collection_id", middleware.AuthorizeJWT, cardHndlr.DontKnowCard)

	router.Run(fmt.Sprintf(":%d", cfg.Port))
}
