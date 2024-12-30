package handler

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	inertia "github.com/romsar/gonertia"
	"neploy.dev/config"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Application struct {
	service service.Application
	inertia *inertia.Inertia
}

func NewApplication(service service.Application, i *inertia.Inertia) *Application {
	return &Application{
		service: service,
		inertia: i,
	}
}

func (a *Application) RegisterRoutes(r *echo.Group) {
	r.POST("", a.Create)
	r.GET("/:id", a.Get)
	r.GET("", a.List)
	r.POST("/:id/deploy", a.Deploy)
	r.POST("/:id/upload", a.Upload)
	r.POST("/:id/start", a.Start)
	r.POST("/:id/stop", a.Stop)
	r.DELETE("/:id", a.Delete)
}

// Create godoc
// @Summary Create a new application
// @Description Create a new application
// @Tags Application
// @Accept json
// @Produce json
// @Param request body model.CreateApplicationRequest true "Application details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications [post]
func (a *Application) Create(c echo.Context) error {
	var req model.CreateApplicationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.AppName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application name is required")
	}

	app := model.Application{
		AppName:     req.AppName,
		Description: req.Description,
	}

	appId, err := a.service.Create(c.Request().Context(), app)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create application")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": appId,
	})
}

// Get godoc
// @Summary Get an application by ID
// @Description Get an application by ID
// @Tags Application
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} model.Application
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications/{id} [get]
func (a *Application) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application ID is required")
	}

	app, err := a.service.Get(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Application not found")
	}

	return c.JSON(http.StatusOK, app)
}

// Deploy godoc
// @Summary Deploy an application
// @Description Deploy an application
// @Tags Application
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param request body model.DeployApplicationRequest true "Deployment details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications/{id}/deploy [post]
func (a *Application) Deploy(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application ID is required")
	}

	var req model.DeployApplicationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	a.service.Deploy(c.Request().Context(), id, req.RepoURL)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "Building",
	})
}

// Start godoc
// @Summary Start an application
// @Description Start an application
// @Tags Application
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications/{id}/start [post]
func (a *Application) Start(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application ID is required")
	}

	err := a.service.StartContainer(c.Request().Context(), id)
	if err != nil {
		logger.Error("error starting application: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start application")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "Running",
	})
}

// Stop godoc
// @Summary Stop an application
// @Description Stop an application
// @Tags Application
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications/{id}/stop [post]
func (a *Application) Stop(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application ID is required")
	}

	if err := a.service.StopContainer(c.Request().Context(), id); err != nil {
		logger.Error("error stopping application: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to stop application")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "Stopped",
	})
}

// Delete godoc
// @Summary Delete an application
// @Description Delete an application
// @Tags Application
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications/{id} [delete]
func (a *Application) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application ID is required")
	}

	if err := a.service.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete application")
	}

	return c.NoContent(http.StatusOK)
}

func (a *Application) Upload(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application ID is required")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File is required")
	}

	path, err := a.service.Upload(c.Request().Context(), id, file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload file")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"storageLocation": filepath.Join(config.Env.UploadPath, path),
	})
}

// List godoc
// @Summary List all applications
// @Description List all applications
// @Tags Application
// @Accept json
// @Produce json
// @Success 200 {object} []model.Application
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications [get]
func (a *Application) List(c echo.Context) error {
	apps, err := a.service.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch applications")
	}

	// If it's a page load (Inertia request), render the full page
	if c.Request().Header.Get("X-Inertia") != "" {
		return a.inertia.Render(c.Response(), c.Request(), "Dashboard/Applications", inertia.Props{
			"applications": apps,
		})
	}

	// For API calls, return JSON
	return c.JSON(http.StatusOK, apps)
}
