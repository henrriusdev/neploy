package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/romsar/gonertia"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Role struct {
	inertia *gonertia.Inertia
	service service.Role
}

func NewRole(i *gonertia.Inertia, service service.Role) *Role {
	return &Role{
		inertia: i,
		service: service,
	}
}

func (h *Role) RegisterRoutes(r *echo.Group) {
	r.GET("", h.List)
	r.POST("", h.Create)
	r.PATCH("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.GET("/users/:id", h.GetUserRoles)
}

// List godoc
// @Summary List all roles
// @Description List all roles
// @Tags Role
// @Accept json
// @Produce json
// @Success 200 {object} []model.Role
// @Failure 500 {object} map[string]interface{}
// @Router /roles [get]
func (h *Role) List(c echo.Context) error {
	roles, err := h.service.Get(c.Request().Context())
	if err != nil {
		logger.Error("error getting roles: %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, roles)
}

// Create godoc
// @Summary Create a new role
// @Description Create a new role
// @Tags Role
// @Accept json
// @Produce json
// @Param request body model.CreateRoleRequest true "Role details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /roles [post]
func (h *Role) Create(c echo.Context) error {
	var req model.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("error binding request: %v", err)
		return h.inertia.Render(c.Response().Writer, c.Request(), "Error/400", nil)
	}

	if err := h.service.Create(c.Request().Context(), req); err != nil {
		logger.Error("error creating role: %v", err)
		return h.inertia.Render(c.Response().Writer, c.Request(), "Error/500", nil)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Role created successfully",
	})
}

// Update godoc
// @Summary Update a role
// @Description Update a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param request body model.CreateRoleRequest true "Role details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /roles/{id} [patch]
func (h *Role) Update(c echo.Context) error {
	id := c.Param("id")
	var req model.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("error binding request: %v", err)
		return h.inertia.Render(c.Response().Writer, c.Request(), "Error/400", nil)
	}

	if err := h.service.Update(c.Request().Context(), id, req); err != nil {
		logger.Error("error updating role: %v", err)
		return h.inertia.Render(c.Response().Writer, c.Request(), "Error/500", nil)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Role updated successfully",
	})
}

// Delete godoc
// @Summary Delete a role
// @Description Delete a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /roles/{id} [delete]
func (h *Role) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		logger.Error("error deleting role: %v", err)
		return h.inertia.Render(c.Response().Writer, c.Request(), "Error/500", nil)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Role deleted successfully",
	})
}

// GetUserRoles godoc
// @Summary Get user roles
// @Description Get roles for a specific user
// @Tags Role
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} []model.UserRoles
// @Failure 500 {object} map[string]interface{}
// @Router /roles/users/{id} [get]
func (h *Role) GetUserRoles(c echo.Context) error {
	userID := c.Param("id")
	roles, err := h.service.GetUserRoles(c.Request().Context(), userID)
	if err != nil {
		logger.Error("error getting user roles: %v", err)
		return h.inertia.Render(c.Response().Writer, c.Request(), "Error/500", nil)
	}

	return h.inertia.Render(c.Response().Writer, c.Request(), "Dashboard/UserRoles", gonertia.Props{
		"roles": roles,
	})
}
