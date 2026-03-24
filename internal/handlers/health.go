package handlers

import (
	"context"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	StatusAlive       = "alive"
	StatusReady       = "ready"
	StatusDegraded    = "degraded"
	StatusUnavailable = "unavailable"

	ServiceName = "stellarbill-backend"

	// Retry Policy Constants
	MaxRetries         = 3
	InitialBackoff     = 500 * time.Millisecond
	MaxDatabaseTimeout = 2 * time.Second
)

// DBPinger defines the minimal interface needed for DB health checks
type DBPinger interface {
	PingContext(ctx context.Context) error
}

// HealthResponse represents the structured health payload
type HealthResponse struct {
	Status       string            `json:"status"`
	Service      string            `json:"service"`
	Timestamp    string            `json:"timestamp"`
	Dependencies map[string]string `json:"dependencies"`
}

// LivenessHandler checks if the service is alive
func LivenessHandler(c *gin.Context) {
	resp := HealthResponse{
		Status:    StatusAlive,
		Service:   ServiceName,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, resp)
}

// ReadinessHandler checks if the service is ready
func ReadinessHandler(db DBPinger) gin.HandlerFunc {
	return func(c *gin.Context) {
		deps := make(map[string]string)

		dbStatus := checkDatabase(db)
		deps["database"] = dbStatus

		overallStatus := deriveOverallStatus(deps)

		resp := HealthResponse{
			Status:       overallStatus,
			Service:      ServiceName,
			Timestamp:    time.Now().UTC().Format(time.RFC3339),
			Dependencies: deps,
		}

		statusCode := http.StatusOK
		if overallStatus == StatusDegraded || overallStatus == StatusUnavailable {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, resp)
	}
}

// checkDatabase implements Timeout and Bounded Retry policies
func checkDatabase(db DBPinger) string {
	if os.Getenv("DATABASE_URL") == "" {
		return "not_configured"
	}
	if db == nil {
		return "down"
	}

	var lastErr error

	// IMPLEMENTATION: Bounded Retry Loop
	for i := 0; i < MaxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), MaxDatabaseTimeout)
		
		lastErr = db.PingContext(ctx)
		cancel() // Release context resources immediately

		if lastErr == nil {
			return "up"
		}

		// If not the last attempt, wait before retrying (Exponential Backoff)
		if i < MaxRetries-1 {
			backoff := time.Duration(math.Pow(2, float64(i))) * InitialBackoff
			time.Sleep(backoff)
		}
	}

	// Determine final failure state
	if lastErr == context.DeadlineExceeded {
		return "timeout"
	}
	return "down"
}

func deriveOverallStatus(deps map[string]string) string {
	for _, status := range deps {
		if status == "down" || status == "timeout" {
			return StatusDegraded
		}
	}
	return StatusReady
}