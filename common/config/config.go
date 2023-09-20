package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI            string
	MongoDatabase       string
	PromotionCollection string
	LogLevel            string
	Port                int
}

// Load carrega as configurações do sistema
func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	return &Config{
		MongoURI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:       getEnv("MONGO_DB", "promotionDB"),
		PromotionCollection: getEnv("PROMOTION_COLLECTION", "promotions"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		Port:                port,
	}, nil
}

// getEnv busca a variável de ambiente
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
