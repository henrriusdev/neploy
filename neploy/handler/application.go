package handler

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"neploy.dev/config"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Application struct {
	service service.Application
}

func NewApplication(service service.Application) *Application {
	return &Application{service: service}
}

func (a *Application) RegisterRoutes(r fiber.Router, i *gonertia.Inertia) {
	r.Post("/", a.Create)
	r.Get("/:id", a.Get)
	r.Post("/:id/deploy", a.Deploy)
	r.Post("/:id/upload", a.Upload)
	r.Post("/:id/start", a.Start)
	r.Post("/:id/stop", a.Stop)
	r.Delete("/:id", a.Delete)
}

type CreateApplicationRequest struct {
	AppName     string `json:"appName"`
	Description string `json:"description"`
	TechStackID string `json:"techStackId"`
}

func (a *Application) Create(c *fiber.Ctx) error {
	var req CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.AppName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application name is required",
		})
	}

	// Create the application
	app := model.Application{
		AppName:     req.AppName,
		Description: req.Description,
		TechStackID: req.TechStackID,
	}

	if err := a.service.Create(c.Context(), app); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create application",
		})
	}

	return c.JSON(fiber.Map{
		"id":      app.ID,
		"message": "Application created successfully",
	})
}

func (a *Application) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application ID is required",
		})
	}

	app, err := a.service.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Application not found",
		})
	}

	return c.JSON(app)
}

func (a *Application) Deploy(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application ID is required",
		})
	}

	var req struct {
		RepoURL string `json:"repoUrl"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Implement deployment logic
	// This should:
	// 1. Clone the repository
	// 2. Detect the tech stack
	// 3. Build the application
	// 4. Deploy it using Docker

	return c.JSON(fiber.Map{
		"message": "Deployment started",
	})
}

func (a *Application) Start(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application ID is required",
		})
	}

	// TODO: Implement start logic
	// This should:
	// 1. Get the application
	// 2. Start the Docker container
	// 3. Update the application status

	return c.JSON(fiber.Map{
		"message": "Application started",
	})
}

func (a *Application) Stop(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application ID is required",
		})
	}

	// TODO: Implement stop logic
	// This should:
	// 1. Get the application
	// 2. Stop the Docker container
	// 3. Update the application status

	return c.JSON(fiber.Map{
		"message": "Application stopped",
	})
}

func (a *Application) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application ID is required",
		})
	}

	// TODO: Implement delete logic
	// This should:
	// 1. Get the application
	// 2. Remove the Docker container and image
	// 3. Delete the application record
	// 4. Clean up any associated files

	return c.JSON(fiber.Map{
		"message": "Application deleted",
	})
}

func (a *Application) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	if filepath.Ext(file.Filename) != ".zip" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file format (only .zip files are allowed)",
		})
	}

	filePath := filepath.Join(config.Env.UploadPath, file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
	})
}
