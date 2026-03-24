package config

import (
	"os"
	"testing"
)

func TestLoadConfigDefaultsAndEnv(t *testing.T) {
	os.Setenv("ENV", "production")
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://example")
	os.Setenv("JWT_SECRET", "jwt")
	os.Setenv("ADMIN_TOKEN", "admintoken")
	os.Setenv("AUDIT_HMAC_SECRET", "secret")
	os.Setenv("AUDIT_LOG_PATH", "custom.log")
	defer func() {
		os.Unsetenv("ENV")
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ADMIN_TOKEN")
		os.Unsetenv("AUDIT_HMAC_SECRET")
		os.Unsetenv("AUDIT_LOG_PATH")
	}()

	cfg := Load()
	if cfg.Env != "production" || cfg.Port != "9090" || cfg.DBConn != "postgres://example" {
		t.Fatalf("env overrides failed: %+v", cfg)
	}
	if cfg.JWTSecret != "jwt" || cfg.AdminToken != "admintoken" {
		t.Fatalf("secrets not loaded: %+v", cfg)
	}
	if cfg.AuditSecret != "secret" || cfg.AuditLogPath != "custom.log" {
		t.Fatalf("audit config not loaded: %+v", cfg)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()
	if cfg.Env != "development" || cfg.Port != "8080" || cfg.DBConn == "" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.AdminToken == "" || cfg.AuditSecret == "" || cfg.AuditLogPath == "" {
		t.Fatalf("expected audit defaults to be present: %+v", cfg)
	}
}
