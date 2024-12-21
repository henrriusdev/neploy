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

func (a *Application) Create(c *fiber.Ctx) error {
	var req model.CreateApplicationRequest
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
	}

	appId, err := a.service.Create(c.Context(), app, req.TechStack)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create application",
		})
	}

	return c.JSON(fiber.Map{
		"id":      appId,
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

	a.service.Deploy(c.Context(), id, req.RepoURL)

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
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application ID is required",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	path, err := a.service.Upload(c.Context(), id, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
		"path":    filepath.Join(config.Env.UploadPath, path),
	})
}
