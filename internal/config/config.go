package config

import "os"

type Config struct {
	Env          string
	Port         string
	DBConn       string
	JWTSecret    string
	AdminToken   string
	AuditSecret  string
	AuditLogPath string
}

func Load() Config {
	return Config{
		Env:          getEnv("ENV", "development"),
		Port:         getEnv("PORT", "8080"),
		DBConn:       getEnv("DATABASE_URL", "postgres://localhost/stellarbill?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "change-me-in-production"),
		AdminToken:   getEnv("ADMIN_TOKEN", "change-me-admin-token"),
		AuditSecret:  getEnv("AUDIT_HMAC_SECRET", "stellarbill-dev-audit"),
		AuditLogPath: getEnv("AUDIT_LOG_PATH", "audit.log"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
