package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestRegister_HealthAndCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	Register(engine)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/health", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected Access-Control-Allow-Origin header to be set")
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if payload["status"] != "ok" {
		t.Fatalf("payload.status = %v, want %q", payload["status"], "ok")
	}
}

func TestRegister_CORSPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	Register(engine)

	req := httptest.NewRequest(http.MethodOptions, "http://localhost:8080/api/health", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatalf("expected Access-Control-Allow-Methods header to be set")
	}
}

func TestRegister_GetSubscriptionRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	Register(engine)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/subscriptions/sub_123", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRegister_GetSubscriptionRejectsInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	Register(engine)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "customer-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("dev-secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/subscriptions/sub_123?expand=plan", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if payload["error"] == nil {
		t.Fatalf("expected payload.error to be present")
	}
}
