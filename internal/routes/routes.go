package routes

import (
	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/audit"
	"stellarbill-backend/internal/config"
	"stellarbill-backend/internal/handlers"
)

func Register(r *gin.Engine, cfg config.Config, auditLogger *audit.Logger) {
	r.Use(corsMiddleware())
	r.Use(audit.Middleware(auditLogger))

	adminHandler := handlers.NewAdminHandler(cfg.AdminToken)

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)
		api.GET("/subscriptions", handlers.ListSubscriptions)
		api.GET("/subscriptions/:id", handlers.GetSubscription)
		api.GET("/plans", handlers.ListPlans)

		admin := api.Group("/admin")
		{
			admin.POST("/purge", adminHandler.PurgeCache)
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
