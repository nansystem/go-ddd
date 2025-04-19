package config

import (
	"os"
	"sync"

	"github.com/nansystem/go-ddd/internal/infrastructure/mysql"
)

type Config struct {
	DBConfig mysql.DBConfig
	GitHub   GitHubConfig
}

var once sync.Once

var config *Config

func LoadConfig() (*Config, error) {
	once.Do(func() {
		config = &Config{}
	})

	dbConfig, err := loadDBConfig()
	if err != nil {
		return nil, err
	}

	config.DBConfig = *dbConfig
	config.GitHub = loadGitHubConfig()

	return config, nil
}

func loadDBConfig() (*mysql.DBConfig, error) {
	dbConfig := mysql.DBConfig{
		User:     getEnv("DB_USER", "ddduser"),
		Password: getEnv("DB_PASSWORD", "dddpass"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "13306"),
		DBName:   getEnv("DB_NAME", "go_ddd"),
	}

	return &dbConfig, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GitHubConfig GitHub設定
type GitHubConfig struct {
	Token string
}

func loadGitHubConfig() GitHubConfig {
	return GitHubConfig{
		Token: getEnv("GITHUB_TOKEN", ""),
	}
}
