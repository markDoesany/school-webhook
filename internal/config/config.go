package config

import "os"

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type FacebookConfig struct {
	VerifyToken     string
	PageAccessToken string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", ""),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", ""),
	}
}

func LoadFacebookConfig() FacebookConfig {
	return FacebookConfig{
		VerifyToken:     getEnv("FB_VERIFY_TOKEN", ""),
		PageAccessToken: getEnv("PAGE_ACCESS_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
