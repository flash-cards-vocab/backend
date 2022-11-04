package application

import (
	"fmt"

	"github.com/flash-cards-vocab/backend/config"
	"github.com/flash-cards-vocab/backend/pkg/database"
)

type Application struct {
	Config    *config.Config
	DBManager *database.Manager
}

func Get() (*Application, error) {
	config := config.New()
	dbManager, err := database.Connect(config)
	if err != nil {
		return nil, err
	}
	fmt.Println("Auto migrate:", config.DBAutoMigrate)
	if config.DBAutoMigrate {
		dbManager.AutoMigrate()
	}

	return &Application{
		Config:    config,
		DBManager: dbManager,
	}, nil
}
