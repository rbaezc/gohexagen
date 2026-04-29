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
	mainContent := ""
	if framework == "Fiber (Recommended - Ultra Fast)" {
		mainContent = `package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Listen(":3000")
}
`
	} else if framework == "Gin (Classic & Battle-tested)" {
		mainContent = `package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})
	r.Run(":3000")
}
`
	} else if framework == "Echo (Elegant & Scalable)" {
		mainContent = `package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.Logger.Fatal(e.Start(":3000"))
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
