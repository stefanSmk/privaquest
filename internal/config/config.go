package config

import "os"

type Config struct {
	Addr        string
	BaseURL     string
	DatabaseURL string
	AdminAPIKey string
	PublicToken string // optional token for public form submissions
}

func Load() Config {
	port := getenv("PORT", "8080")
	return Config{
		Addr:        ":" + port,
		BaseURL:     getenv("BASE_URL", "http://localhost:"+port),
		DatabaseURL: getenv("DATABASE_URL", "file:privaquest.db?_pragma=foreign_keys(1)"),
		AdminAPIKey: getenv("ADMIN_API_KEY", "change-me-admin-key"),
		PublicToken: getenv("PUBLIC_TOKEN", ""),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
