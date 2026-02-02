package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Port        string
	DataDir     string
	DatabasePath string
	MediaDir    string
	BaseURL     string
}

func Load() *Config {
	dataDir := getEnv("DATA_DIR", "./data")

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DataDir:     dataDir,
		DatabasePath: filepath.Join(dataDir, "listenbucket.db"),
		MediaDir:    filepath.Join(dataDir, "media"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
