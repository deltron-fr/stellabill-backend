package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/audit"
	"stellarbill-backend/internal/config"
	"stellarbill-backend/internal/routes"
)

func main() {
	if err := mainWithRunner(nil); err != nil {
		log.Fatal(err)
	}
}

// mainWithRunner allows tests to inject a runner while covering the main entrypoint.
func mainWithRunner(runner func(*gin.Engine, string) error) error {
	cfg := config.Load()
	if runner == nil {
		runner = func(r *gin.Engine, addr string) error { return r.Run(addr) }
	}
	if os.Getenv("SKIP_SERVER_RUN") == "1" {
		runner = func(_ *gin.Engine, _ string) error { return nil }
	}
	return run(cfg, runner)
}

// buildRouter is split from main to make testing and audit wiring easier.
func buildRouter(cfg config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	auditLogger := audit.NewLogger(cfg.AuditSecret, audit.NewFileSink(cfg.AuditLogPath))
	routes.Register(router, cfg, auditLogger)
	return router
}

// run builds the router and delegates serving, allowing tests to stub the runner.
func run(cfg config.Config, runner func(*gin.Engine, string) error) error {
	router := buildRouter(cfg)
	addr := ":" + cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("Stellarbill backend listening on %s", addr)
	return runner(router, addr)
}
