package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/config"
)

func TestBuildRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		Env:         "test",
		Port:        "8080",
		AdminToken:  "token",
		AuditSecret: "secret",
	}
	router := buildRouter(cfg)

	req, _ := http.NewRequest("GET", "/api/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("health endpoint not registered, got %d", rec.Code)
	}
}

func TestBuildRouterSetsReleaseModeForProduction(t *testing.T) {
	prevMode := gin.Mode()
	defer gin.SetMode(prevMode)
	cfg := config.Config{
		Env:         "production",
		Port:        "8080",
		AdminToken:  "token",
		AuditSecret: "secret",
	}
	router := buildRouter(cfg)
	if gin.Mode() != gin.ReleaseMode {
		t.Fatalf("expected release mode for production env, got %s", gin.Mode())
	}
	if router == nil {
		t.Fatal("router should not be nil")
	}
}

func TestRunRespectsPortOverride(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		Env:         "test",
		Port:        "8080",
		AdminToken:  "token",
		AuditSecret: "secret",
	}
	t.Setenv("PORT", "9090")
	var capturedAddr string
	err := run(cfg, func(_ *gin.Engine, addr string) error {
		capturedAddr = addr
		return nil
	})
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if capturedAddr != ":9090" {
		t.Fatalf("expected addr :9090, got %s", capturedAddr)
	}
	os.Unsetenv("PORT")
}

func TestMainSkipsWhenFlagged(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("SKIP_SERVER_RUN", "1")
	t.Setenv("AUDIT_LOG_PATH", filepath.Join(t.TempDir(), "audit.log"))
	main()
}

func TestMainWithRunnerUsesInjection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("SKIP_SERVER_RUN", "")
	called := false
	err := mainWithRunner(func(_ *gin.Engine, addr string) error {
		called = true
		if addr == "" {
			t.Fatal("address should not be empty")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected injected runner to be invoked")
	}
}
