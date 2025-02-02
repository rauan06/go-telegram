package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBPort     string `env:"DB_PORT"`
}

var (
	cfg Config
)

func LoadConfig() *Config {
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBUser = getEnv("DB_USER", "postgres")
	cfg.DBPassword = getEnv("DB_PASSWORD", "0000")
	cfg.DBName = getEnv("DB_NAME", "frappuccino_db")
	cfg.DBPort = getEnv("DB_PORT", "5432")

	return &cfg
}

func GetConfing() *Config {
	return &cfg
}

func (c *Config) MakeConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
