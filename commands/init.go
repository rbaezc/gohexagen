package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitProject(projectName, framework, orm string) error {
	basePath := projectName
	dirs := []string{
		"cmd",
		"internal/domain/models",
		"internal/domain/repositories",
		"internal/application/usecases",
		"internal/infrastructure/adapters/http",
		"internal/infrastructure/adapters/db",
		"internal/infrastructure/config",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(basePath, dir), 0755); err != nil {
			return err
		}
	}

	// Write go.mod
	goMod := fmt.Sprintf("module %s\n\ngo 1.21\n", projectName)
	if err := os.WriteFile(filepath.Join(basePath, "go.mod"), []byte(goMod), 0644); err != nil {
		return err
	}

	// Write cmd/main.go
	// Write cmd/main.go
	mainContent := ""
	if framework == "Fiber (Recommended - Ultra Fast)" {
		mainContent = `package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	go func() {
		slog.Info("Starting Fiber HTTP server on :3000")
		if err := app.Listen(":3000"); err != nil {
			slog.Error("Server crash", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server elegantly...")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		slog.Error("Shutdown failed", "err", err)
	}
	slog.Info("Server gracefully stopped")
}
`
	} else if framework == "Gin (Classic & Battle-tested)" {
		mainContent = `package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	srv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	go func() {
		slog.Info("Starting Gin HTTP server on :3000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server crash", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server elegantly...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Shutdown failed", "err", err)
	}
	slog.Info("Server gracefully stopped")
}
`
	} else if framework == "Echo (Elegant & Scalable)" {
		mainContent = `package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	e := echo.New()
	e.HideBanner = true

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	go func() {
		slog.Info("Starting Echo HTTP server on :3000")
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			slog.Error("Server crash", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server elegantly...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		slog.Error("Shutdown failed", "err", err)
	}
	slog.Info("Server gracefully stopped")
}
`
	} else {
		mainContent = `package main

import "fmt"

func main() {
	fmt.Println("Hello Hexagonal Architecture in Go!")
}
`
	}

	if err := os.WriteFile(filepath.Join(basePath, "cmd/main.go"), []byte(mainContent), 0644); err != nil {
		return err
	}

	return nil
}
