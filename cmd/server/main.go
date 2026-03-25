package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/config"
	"stellarbill-backend/internal/routes"
)

var listenAndServe = func(srv *http.Server) error {
	return srv.ListenAndServe()
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		printConfigError(err)
		os.Exit(1)
	}

	if vResult := cfg.Validate(); len(vResult.Warnings) > 0 {
		log.Printf("WARNING: Configuration warnings:")
		for _, warning := range vResult.Warnings {
			log.Printf("  - %s", warning)
		}
	}

	if err := run(cfg); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func printConfigError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: Configuration validation failed: %s\n", err.Error())
	fmt.Fprintln(os.Stderr, "\nRequired environment variables:")
	for _, key := range config.GetRequiredEnvVars() {
		fmt.Fprintf(os.Stderr, "  - %s\n", key)
	}
	fmt.Fprintln(os.Stderr, "\nOptional environment variables and defaults:")
	for key, val := range config.GetOptionalEnvVars() {
		fmt.Fprintf(os.Stderr, "  - %s (default: %s)\n", key, val)
	}
}

func run(cfg config.Config) error {
	configureGinMode(cfg.Env)

	srv := newServer(cfg)
	log.Printf("Starting Stellarbill backend on %s (env: %s)", srv.Addr, cfg.Env)
	log.Printf("Server timeouts - Read: %ds, Write: %ds, Idle: %ds", cfg.ReadTimeout, cfg.WriteTimeout, cfg.IdleTimeout)

	return listenAndServe(srv)
}

func configureGinMode(env string) {
	switch env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
		log.Println("Running in production mode")
	case "development":
		gin.SetMode(gin.DebugMode)
		log.Println("Running in development mode")
	default:
		gin.SetMode(gin.TestMode)
		log.Printf("Running in %s mode", env)
	}
}

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})
	routes.Register(router)
	return router
}

func newServer(cfg config.Config) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        newRouter(),
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(cfg.IdleTimeout) * time.Second,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
}
