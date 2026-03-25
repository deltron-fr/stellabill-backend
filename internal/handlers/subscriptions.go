package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"stellarbill-backend/internal/requestparams"
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
	if _, err := requestparams.SanitizeQuery(c.Request.URL.Query(), requestparams.QueryRules{
		Strings: map[string]requestparams.StringRule{
			"customer": requestparams.IdentifierRule(64),
			"plan_id":  requestparams.IdentifierRule(64),
			"status":   requestparams.EnumRule(16, true, "active", "past_due", "canceled", "trialing"),
		},
		Ints: map[string]requestparams.IntRule{
			"limit": {Min: 1, Max: 100},
			"page":  {Min: 1, Max: 100000},
		},
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: load from DB, filter by merchant from JWT/API key
	subscriptions := []Subscription{}
	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

// NewGetSubscriptionHandler returns a gin.HandlerFunc that retrieves a full
// subscription detail using the provided SubscriptionService.
func NewGetSubscriptionHandler(svc service.SubscriptionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		callerID, exists := c.Get("callerID")
		if !exists {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if _, err := requestparams.SanitizeQuery(c.Request.URL.Query(), requestparams.QueryRules{}); err != nil {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := requestparams.NormalizePathID("id", c.Param("id"))
		if err != nil {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		detail, warnings, err := svc.GetDetail(c.Request.Context(), callerID.(string), id)
		if err != nil {
			c.Header("Content-Type", "application/json; charset=utf-8")
			switch err {
			case service.ErrNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			case service.ErrDeleted:
				c.JSON(http.StatusGone, gin.H{"error": "subscription has been deleted"})
			case service.ErrForbidden:
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			case service.ErrBillingParse:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			}
			return
		}

		c.Header("Content-Type", "application/json; charset=utf-8")
		envelope := service.ResponseEnvelope{
			APIVersion: "1",
			Data:       detail,
			Warnings:   warnings,
		}
		c.JSON(http.StatusOK, envelope)
	}
}
