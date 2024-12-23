package handler

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"neploy.dev/config"
	"neploy.dev/pkg/logger"
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
	r.Post("", a.Create)
	r.Get("/:id", a.Get)
	r.Get("", a.List)
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

	res, err := a.service.StartContainer(c.Context(), id)
	if err != nil {
		logger.Error("error starting application: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start application",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Application started",
		"containerId": res,
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

	if err := a.service.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete application",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Application deleted successfully",
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

func (a *Application) List(c *fiber.Ctx) error {
	apps, err := a.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list applications",
		})
	}

	return c.JSON(apps)
}
