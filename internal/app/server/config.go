package server

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}
	CryptoAPI struct {
		Token string
	}
	WorkerPool struct {
		Size       int
		UpdateTime int
	}
}

func LoadConfig() *Config {
	var cfg Config

	// Server
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// Database
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.Name = getEnv("DB_NAME", "crypto")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")

	// CryptoAPI
	cfg.CryptoAPI.Token = getEnv("CRYPTO_API_KEY", "")

	// WorkerPool
	cfg.WorkerPool.Size, _ = strconv.Atoi(getEnv("WORKER_POOL_SIZE", "10"))
	cfg.WorkerPool.UpdateTime, _ = strconv.Atoi(getEnv("WORKER_POOL_UPDATE_TIME", "60"))

	// Validate
	if cfg.Database.Password == "" {
		log.Fatal("DB_PASSWORD is required")
	}
	if cfg.WorkerPool.Size == 0 || cfg.WorkerPool.UpdateTime == 0 {
		log.Fatal("WORKER_POOL_SIZE and WORKER_POOL_UPDATE_TIME must be int and greater than 0")
	}
	return &cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDBConnectionString возвращает строку подключения к PostgreSQL
func (c *Config) GetDBConnectionString() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" dbname=" + c.Database.Name +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" sslmode=disable"
}
