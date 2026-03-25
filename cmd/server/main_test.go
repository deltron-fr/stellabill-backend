package main

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"stellarbill-backend/internal/config"
)

func TestConfigureGinMode(t *testing.T) {
	configureGinMode("production")
	if gin.Mode() != gin.ReleaseMode {
		t.Fatalf("gin.Mode = %q, want %q", gin.Mode(), gin.ReleaseMode)
	}

	configureGinMode("development")
	if gin.Mode() != gin.DebugMode {
		t.Fatalf("gin.Mode = %q, want %q", gin.Mode(), gin.DebugMode)
	}

	configureGinMode("staging")
	if gin.Mode() != gin.TestMode {
		t.Fatalf("gin.Mode = %q, want %q", gin.Mode(), gin.TestMode)
	}
}

func TestNewRouterRegistersExpectedRoutes(t *testing.T) {
	router := newRouter()
	if len(router.Routes()) == 0 {
		t.Fatal("expected router to register routes")
	}
}

func TestNewServerUsesConfiguredSettings(t *testing.T) {
	cfg := config.Config{
		Env:            "development",
		Port:           9191,
		MaxHeaderBytes: 2048,
		ReadTimeout:    10,
		WriteTimeout:   20,
		IdleTimeout:    30,
	}

	srv := newServer(cfg)

	if srv.Addr != ":9191" {
		t.Fatalf("Addr = %q, want %q", srv.Addr, ":9191")
	}
	if srv.MaxHeaderBytes != 2048 {
		t.Fatalf("MaxHeaderBytes = %d, want 2048", srv.MaxHeaderBytes)
	}
	if srv.ReadTimeout != 10*time.Second {
		t.Fatalf("ReadTimeout = %v, want %v", srv.ReadTimeout, 10*time.Second)
	}
	if srv.WriteTimeout != 20*time.Second {
		t.Fatalf("WriteTimeout = %v, want %v", srv.WriteTimeout, 20*time.Second)
	}
	if srv.IdleTimeout != 30*time.Second {
		t.Fatalf("IdleTimeout = %v, want %v", srv.IdleTimeout, 30*time.Second)
	}
}

func TestRunUsesConfiguredServer(t *testing.T) {
	previous := listenAndServe
	defer func() { listenAndServe = previous }()

	called := false
	listenAndServe = func(srv *http.Server) error {
		called = true
		if srv.Addr != ":9191" {
			t.Fatalf("Addr = %q, want %q", srv.Addr, ":9191")
		}
		if srv.Handler == nil {
			t.Fatal("expected server handler to be initialized")
		}
		return nil
	}

	err := run(config.Config{
		Env:            "production",
		Port:           9191,
		MaxHeaderBytes: 1024,
		ReadTimeout:    10,
		WriteTimeout:   10,
		IdleTimeout:    10,
	})
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !called {
		t.Fatal("expected listenAndServe to be called")
	}
	if gin.Mode() != gin.ReleaseMode {
		t.Fatalf("gin.Mode = %q, want %q", gin.Mode(), gin.ReleaseMode)
	}
}

func TestRunPropagatesServerError(t *testing.T) {
	previous := listenAndServe
	defer func() { listenAndServe = previous }()

	wantErr := errors.New("listen failed")
	listenAndServe = func(_ *http.Server) error {
		return wantErr
	}

	err := run(config.Config{
		Env:            "development",
		Port:           8080,
		MaxHeaderBytes: 1024,
		ReadTimeout:    10,
		WriteTimeout:   10,
		IdleTimeout:    10,
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("run error = %v, want %v", err, wantErr)
	}
}
