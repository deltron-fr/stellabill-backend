package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"stellarbill-backend/internal/service"
)

type Subscription struct {
	ID          string `json:"id"`
	PlanID      string `json:"plan_id"`
	Customer    string `json:"customer"`
	Status      string `json:"status"`
	Amount      string `json:"amount"`
	Interval    string `json:"interval"`
	NextBilling string `json:"next_billing,omitempty"`
}

func ListSubscriptions(c *gin.Context) {
	// TODO: load from DB, filter by merchant from JWT/API key
	subscriptions := []Subscription{}
	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

// NewGetSubscriptionHandler returns a gin.HandlerFunc that retrieves a full
// subscription detail using the provided SubscriptionService.
func NewGetSubscriptionHandler(svc service.SubscriptionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Read callerID and tenantID from context (set by AuthMiddleware).
		callerID, callerExists := c.Get("callerID")
		tenantID, tenantExists := c.Get("tenantID")
		if !callerExists || !tenantExists {
			RespondWithAuthError(c, "caller or tenant information missing from context")
			return
		}

		// 2. Validate :id path param.
		id := c.Param("id")
		if strings.TrimSpace(id) == "" {
			RespondWithValidationError(c, "subscription id is required", map[string]interface{}{
				"field": "id",
				"reason": "cannot be empty",
			})
			return
		}

		// 3. Call service with tenant scope.
		detail, warnings, err := svc.GetDetail(c.Request.Context(), tenantID.(string), callerID.(string), id)
		if err != nil {
			statusCode, errorCode, message := MapServiceErrorToResponse(err)
			RespondWithError(c, statusCode, errorCode, message)
			return
		}

		// 4. Set Content-Type and respond with envelope.
		c.Header("Content-Type", "application/json; charset=utf-8")
		envelope := service.ResponseEnvelope{
			APIVersion: "1",
			Data:       detail,
			Warnings:   warnings,
		}
		c.JSON(http.StatusOK, envelope)
	}
}
