package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/health", Health)
	r.GET("/api/plans", ListPlans)
	r.GET("/api/subscriptions", ListSubscriptions)
	return r
}

func TestListPlans(t *testing.T) {
	r := testRouter()

	t.Run("accepts normalized filter inputs", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/plans?currency=ngn&interval=%20MONTHLY%20&search=Starter%20Plan&limit=%EF%BC%91%EF%BC%90&page=2", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body["plans"]; !ok {
			t.Fatalf("expected plans key in response body")
		}
	})

	t.Run("rejects malicious search payloads", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/plans?search=%3Cscript%3E", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		assertErrorContains(t, rec, "invalid query parameter \"search\"")
	})
}

func TestHealth(t *testing.T) {
	r := testRouter()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status = %q, want %q", body["status"], "ok")
	}
}

func TestListSubscriptions(t *testing.T) {
	r := testRouter()

	t.Run("rejects unsupported query parameters", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/subscriptions?debug=true", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		assertErrorContains(t, rec, "unsupported parameter")
	})

	t.Run("rejects overflow pagination values", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/subscriptions?page=999999999999999999999999", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		assertErrorContains(t, rec, "valid integer")
	})

	t.Run("accepts normalized status filters", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/subscriptions?status=%20ACTIVE%20&limit=25&page=%EF%BC%91", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}

func assertErrorContains(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got := body["error"]; got == "" || !strings.Contains(got, want) {
		t.Fatalf("error = %q, want substring %q", got, want)
	}
}
