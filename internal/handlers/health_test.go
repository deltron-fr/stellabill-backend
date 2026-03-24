package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
)

// --------------------
// MOCK DB
// --------------------

type MockDB struct {
	ShouldFail    bool
	ShouldTimeout bool
	CallCount     int32 // Tracks how many times PingContext was called
}

func (m *MockDB) PingContext(ctx context.Context) error {
	atomic.AddInt32(&m.CallCount, 1) // Increment count safely

	if m.ShouldTimeout {
		// Mock a delay that exceeds the 2s context timeout
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if m.ShouldFail {
		return errors.New("db failure")
	}
	return nil
}

// --------------------
// TESTS
// --------------------

func TestLivenessHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/live", LivenessHandler)

	req, _ := http.NewRequest("GET", "/live", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestReadiness_Healthy(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test")
	db := &MockDB{}

	router := gin.Default()
	router.GET("/ready", ReadinessHandler(db))

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if db.CallCount != 1 {
		t.Errorf("expected 1 call, got %d", db.CallCount)
	}
}

func TestReadiness_DBRetryLogic(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test")
	
	// Case: Database fails, should retry 3 times
	db := &MockDB{ShouldFail: true}

	router := gin.Default()
	router.GET("/ready", ReadinessHandler(db))

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}

	// CRITICAL: This proves our MaxRetries logic is working!
	if db.CallCount != 3 {
		t.Errorf("expected 3 retry attempts, got %d", db.CallCount)
	}
}

func TestReadiness_DBTimeoutWithRetry(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test")
	db := &MockDB{ShouldTimeout: true}

	router := gin.Default()
	router.GET("/ready", ReadinessHandler(db))

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
	
	// Ensure it tried multiple times before timing out completely
	if db.CallCount != 3 {
		t.Errorf("expected 3 timeout attempts, got %d", db.CallCount)
	}
}

func TestReadiness_NoDBConfigured(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	router := gin.Default()
	router.GET("/ready", ReadinessHandler(nil))

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}