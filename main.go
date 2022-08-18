package main

import (
	"fmt"
	"log"

	collection_uc "github.com/flash-cards-vocab/backend/app/usecase/collection"
	user_uc "github.com/flash-cards-vocab/backend/app/usecase/user"
	"github.com/flash-cards-vocab/backend/config"
	"github.com/flash-cards-vocab/backend/handler"
	"github.com/flash-cards-vocab/backend/middleware"
	"github.com/flash-cards-vocab/backend/repository/collection_repository"
	"github.com/flash-cards-vocab/backend/repository/user_repository"
	"github.com/gin-gonic/gin"
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
	collection_uc := collection_uc.New(collection_repo)

	collectionHndlr := handler.NewCollectionHandler(collection_uc)

	/* CRUD author for CMS */
	collection := router.Group("/collection")
	collection.GET("/my", middleware.AuthorizeJWT, collectionHndlr.GetMyCollections)
	collection.PUT("/like/:id", middleware.AuthorizeJWT, collectionHndlr.LikeCollectionById)
	collection.PUT("/dislike/:id", middleware.AuthorizeJWT, collectionHndlr.DislikeCollectionById)
	collection.PUT("/view/:id", middleware.AuthorizeJWT, collectionHndlr.ViewCollectionById)

	// card := router.Group("/card")
	// card.GET("/get-first-ten", h.GetMyCollections)
	// card.PUT("/register", h.LikeCollectionById)
	// card.PUT("/logout", h.DislikeCollectionById)
	// card.PUT("/refresh-token", h.ViewCollectionById)

	router.Run(fmt.Sprintf(":%d", cfg.Port))
}
