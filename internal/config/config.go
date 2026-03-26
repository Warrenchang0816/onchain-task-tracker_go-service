package config

import "os"

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
}

func Load() Config {
	return Config{
		AppPort:   getEnv("APP_PORT", ""),
		DBHost:    getEnv("DB_HOST", ""),
		DBPort:    getEnv("DB_PORT", ""),
		DBUser:    getEnv("DB_USER", ""),
		DBPass:    getEnv("DB_PASS", ""),
		DBName:    getEnv("DB_NAME", ""),
		DBSSLMode: getEnv("DB_SSLMODE", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
