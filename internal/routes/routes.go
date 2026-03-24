package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/auth"
	"stellarbill-backend/internal/config"
	"stellarbill-backend/internal/handlers"
	"stellarbill-backend/internal/middleware"
	"stellarbill-backend/internal/repository"
	"stellarbill-backend/internal/service"
)

func Register(r *gin.Engine) {
	cfg := config.Load()

	r.Use(corsMiddleware())

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}

	subRepo := repository.NewMockSubscriptionRepo()
	planRepo := repository.NewMockPlanRepo()
	svc := service.NewSubscriptionService(subRepo, planRepo)

	api := r.Group("/api")
	{
		// Public health check - no authentication required
		api.GET("/health", handlers.Health)

		// Authenticated routes requiring authentication
		authenticated := api.Group("")
		authenticated.Use(auth.AuthMiddleware(cfg.JWTSecret))
		{
			// List plans - any authenticated user
			authenticated.GET("/plans", handlers.ListPlans)

			// Subscriptions - requires admin or merchant role
			authenticated.GET("/subscriptions", auth.AuthzMiddleware(auth.RoleAdmin, auth.RoleMerchant), handlers.ListSubscriptions)
			authenticated.GET("/subscriptions/:id", auth.AuthzMiddleware(auth.RoleAdmin, auth.RoleMerchant), handlers.GetSubscription)
		}

	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
