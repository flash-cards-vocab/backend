package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"rest-api"`
	Environment string `envconfig:"ENV" default:"dev"`
	Port        int    `envconfig:"PORT" default:"8080" required:"true"`

	DBHost         string `envconfig:"DB_HOST" default:"localhost"`
	DBPort         string `envconfig:"DB_PORT" default:"33062"`
	DBUserName     string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword     string `envconfig:"DB_PASSWORD" default:"password"`
	DBDatabaseName string `envconfig:"DB_DBNAME" default:"gca"`
	DBLogMode      int    `envconfig:"DB_LOG_MODE" default:"3"`

	RedisHost string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort string `envconfig:"REDIS_PORT" default:"33792"`

	GCSBucketName string `envconfig:"GCS_BUCKET_NAME" default:"flashcards-images"`
	GCSPrefix     string `envconfig:"GCS_PREFIX" default:"dev"`
	GCSAPIKey     string `envconfig:"GCS_API_KEY" default:""`
	GCSJSONAPIKey string `envconfig:"GCS_JSON_API_KEY" default:""`
}

func New() *Config {
	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return &config

}
