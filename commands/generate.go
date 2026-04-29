package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "project"
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 && parts[0] == "module" {
			return parts[1]
		}
	}
	return "project"
}

func toTitle(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func GenerateResource(entityName string, fields []string) error {
	moduleName := getModuleName()
	
	// Standardize casing
	structName := toTitle(entityName)
	snakeName := strings.ToLower(entityName)

	// Parse fields: format name:type
	var structFields []string
	for _, f := range fields {
		parts := strings.Split(f, ":")
		if len(parts) == 2 {
			goType := "string"
			switch parts[1] {
			case "string":
				goType = "string"
			case "number":
				goType = "int"
			case "float":
				goType = "float64"
			case "boolean":
				goType = "bool"
			}
			fieldName := toTitle(parts[0])
			structFields = append(structFields, fmt.Sprintf("\t%s %s `json:\"%s\"`", fieldName, goType, parts[0]))
		}
	}

	// 1. Domain Model
	modelContent := fmt.Sprintf("package models\n\ntype %s struct {\n\tID string `json:\"id\"` \n%s\n}\n", structName, strings.Join(structFields, "\n"))

	modelPath := filepath.Join("internal/domain/models", snakeName+".go")
	_ = os.MkdirAll(filepath.Dir(modelPath), 0755)
	if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
		return err
	}

	// 2. Domain Repository Interface (Port)
	repoContent := fmt.Sprintf(`package repositories

import "%s/internal/domain/models"

type I%sRepository interface {
	FindByID(id string) (*models.%s, error)
	Save(entity *models.%s) error
}
`, moduleName, structName, structName, structName)

	repoPath := filepath.Join("internal/domain/repositories", snakeName+"_repository.go")
	_ = os.MkdirAll(filepath.Dir(repoPath), 0755)
	if err := os.WriteFile(repoPath, []byte(repoContent), 0644); err != nil {
		return err
	}

	// 3. Application UseCase
	usecaseContent := fmt.Sprintf(`package usecases

import (
	"%s/internal/domain/models"
	"%s/internal/domain/repositories"
)

type Create%sUseCase struct {
	Repo repositories.I%sRepository
}

func NewCreate%sUseCase(repo repositories.I%sRepository) *Create%sUseCase {
	return &Create%sUseCase{Repo: repo}
}

func (u *Create%sUseCase) Execute(entity *models.%s) error {
	return u.Repo.Save(entity)
}
`, moduleName, moduleName, structName, structName, structName, structName, structName, structName, structName, structName)

	usecasePath := filepath.Join("internal/application/usecases", snakeName+".go")
	_ = os.MkdirAll(filepath.Dir(usecasePath), 0755)
	if err := os.WriteFile(usecasePath, []byte(usecaseContent), 0644); err != nil {
		return err
	}

	// 4. Infrastructure HTTP Adapter
	mainData, _ := os.ReadFile("cmd/main.go")
	mainStr := string(mainData)

	framework := "Fiber"
	if strings.Contains(mainStr, "gin") {
		framework = "Gin"
	} else if strings.Contains(mainStr, "echo") {
		framework = "Echo"
	}

	var httpContent string
	switch framework {
	case "Fiber":
		httpContent = fmt.Sprintf(`package http

import (
	"github.com/gofiber/fiber/v2"
	"%s/internal/application/usecases"
	"%s/internal/domain/models"
)

type %sHandler struct {
	UseCase *usecases.Create%sUseCase
}

func New%sHandler(uc *usecases.Create%sUseCase) *%sHandler {
	return &%sHandler{UseCase: uc}
}

func (h *%sHandler) Create(c *fiber.Ctx) error {
	var entity models.%s
	if err := c.BodyParser(&entity); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.UseCase.Execute(&entity); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(entity)
}
`, moduleName, moduleName, structName, structName, structName, structName, structName, structName, structName, structName)
	case "Gin":
		httpContent = fmt.Sprintf(`package http

import (
	"github.com/gin-gonic/gin"
	"%s/internal/application/usecases"
	"%s/internal/domain/models"
)

type %sHandler struct {
	UseCase *usecases.Create%sUseCase
}

func New%sHandler(uc *usecases.Create%sUseCase) *%sHandler {
	return &%sHandler{UseCase: uc}
}

func (h *%sHandler) Create(c *gin.Context) {
	var entity models.%s
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.UseCase.Execute(&entity); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, entity)
}
`, moduleName, moduleName, structName, structName, structName, structName, structName, structName, structName, structName)
	default:
		httpContent = fmt.Sprintf(`package http

import (
	"github.com/labstack/echo/v4"
	"%s/internal/application/usecases"
	"%s/internal/domain/models"
)

type %sHandler struct {
	UseCase *usecases.Create%sUseCase
}

func New%sHandler(uc *usecases.Create%sUseCase) *%sHandler {
	return &%sHandler{UseCase: uc}
}

func (h *%sHandler) Create(c echo.Context) error {
	var entity models.%s
	if err := c.Bind(&entity); err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	if err := h.UseCase.Execute(&entity); err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(201, entity)
}
`, moduleName, moduleName, structName, structName, structName, structName, structName, structName, structName, structName)
	}

	httpPath := filepath.Join("internal/infrastructure/adapters/http", snakeName+"_handler.go")
	_ = os.MkdirAll(filepath.Dir(httpPath), 0755)
	_ = os.WriteFile(httpPath, []byte(httpContent), 0644)

	fmt.Printf("Successfully generated domain models, ports, usecases, and HTTP handlers for %s!\n", structName)
	return nil
}
