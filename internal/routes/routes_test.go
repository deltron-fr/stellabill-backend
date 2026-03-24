package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/audit"
	"stellarbill-backend/internal/config"
)

func TestRoutesRegistrationAndCors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		Port:        "8080",
		AdminToken:  "token",
		AuditSecret: "secret",
	}
	sink := &audit.MemorySink{}
	logger := audit.NewLogger(cfg.AuditSecret, sink)

	r := gin.New()
	Register(r, cfg, logger)

	req, _ := http.NewRequest("OPTIONS", "/api/subscriptions", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for preflight, got %d", rec.Code)
	}
	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Fatalf("cors header missing")
	}

	// Regular request passes through middleware chain
	getReq, _ := http.NewRequest("GET", "/api/health", nil)
	getRec := httptest.NewRecorder()
	r.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for health, got %d", getRec.Code)
	}
	if getRec.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Fatalf("cors headers not attached on normal request")
	}
}
