package database

import (
	"fmt"
	"time"

	"github.com/flash-cards-vocab/backend/config"
	cardRepo "github.com/flash-cards-vocab/backend/pkg/repository/card_repository"
	collectionRepo "github.com/flash-cards-vocab/backend/pkg/repository/collection_repository"
	userRepo "github.com/flash-cards-vocab/backend/pkg/repository/user_repository"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type postgres struct {
	DBHost                        string
	DBPort                        string
	DBUserName                    string
	DBPass                        string
	DBDatabaseName                string
	DBLogMode                     logger.LogLevel
	maxIdleConnection             int
	maxOpenConnection             int
	connectionMaxLifetimeInSecond int
}
type mysqlOption func(*postgres)

type Manager struct {
	DB *gorm.DB
}

func Connect(config *config.Config) (*Manager, error) {
	db := &postgres{
		DBHost:                        config.DBHost,
		DBPort:                        config.DBPort,
		DBUserName:                    config.DBUserName,
		DBPass:                        config.DBPassword,
		DBDatabaseName:                config.DBDatabaseName,
		DBLogMode:                     logger.LogLevel(config.DBLogMode),
		maxIdleConnection:             5,
		maxOpenConnection:             10,
		connectionMaxLifetimeInSecond: 60,
	}

	// for _, o := range config.DBOptions {
	// 	o(db)
	// }

	return connect(db)
}

func (m *Manager) AutoMigrate() {
	fmt.Println("auto migrating...")
	m.DB.AutoMigrate(
		cardRepo.Card{},
		cardRepo.CardMetrics{},
		cardRepo.CardUserProgress{},
		cardRepo.CollectionCards{},
		cardRepo.CollectionUserProgress{},
		collectionRepo.Collection{},
		collectionRepo.CollectionCards{},
		collectionRepo.CollectionMetrics{},
		collectionRepo.CollectionUserMetrics{},
		collectionRepo.CollectionUserProgress{},
		userRepo.User{},
	)
}

func connect(param *postgres) (*Manager, error) {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=Local&parseTime=true",
	// 	param.DBUserName, param.DBPass, param.DBHost, param.DBPort, param.DBDatabaseName)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", //required",
		param.DBHost, param.DBUserName, param.DBPass, param.DBDatabaseName, param.DBPort)

	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(param.DBLogMode)),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	// set configuration pooling connection
	mysqlDb, _ := db.DB()
	mysqlDb.SetMaxOpenConns(param.maxOpenConnection)
	mysqlDb.SetConnMaxLifetime(time.Duration(param.connectionMaxLifetimeInSecond) * time.Minute)
	mysqlDb.SetMaxIdleConns(param.maxIdleConnection)

	return &Manager{
		DB: db,
	}, nil
}

func SetMaxIdleConns(conns int) mysqlOption {
	return func(c *postgres) {
		if conns > 0 {
			c.maxIdleConnection = conns
		}
	}
}

func SetMaxOpenConns(conns int) mysqlOption {
	return func(c *postgres) {
		if conns > 0 {
			c.maxOpenConnection = conns
		}
	}
}

func SetConnMaxLifetime(conns int) mysqlOption {
	return func(c *postgres) {
		if conns > 0 {
			c.connectionMaxLifetimeInSecond = conns
		}
	}
}
