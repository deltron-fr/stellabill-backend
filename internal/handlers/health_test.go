package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Health(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestSubscriptionHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// List
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ListSubscriptions(c)
		if w.Code != http.StatusOK {
			t.Fatalf("list status %d", w.Code)
		}
	}

	// Get missing id
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{}
		GetSubscription(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for missing id, got %d", w.Code)
		}
	}

	// Get with id
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "sub_123"}}
		GetSubscription(c)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 for valid id, got %d", w.Code)
		}
	}
}

func TestListPlans(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ListPlans(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
