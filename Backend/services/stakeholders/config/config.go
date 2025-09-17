package config

import "os"

type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBName string
	DBUser string
	DBPass string
}

func NewConfig() *Config {
	return &Config{
		Port:   os.Getenv("USER_SERVICE_PORT"),
		DBHost: os.Getenv("USER_DB_HOST"),
		DBPort: os.Getenv("USER_DB_PORT"),
		DBName: os.Getenv("USER_DB_NAME"),
		DBUser: os.Getenv("USER_DB_USER"),
		DBPass: os.Getenv("USER_DB_PASS"),
	}
}
