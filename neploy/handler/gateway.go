package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inertia "github.com/romsar/gonertia"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Gateway struct {
	gatewayService service.Gateway
	healthChecker  service.HealthChecker
}

func NewGateway(gatewayService service.Gateway, healthChecker service.HealthChecker) *Gateway {
	return &Gateway{
		gatewayService: gatewayService,
		healthChecker:  healthChecker,
	}
}

func (h *Gateway) RegisterRoutes(r *echo.Group, i *inertia.Inertia) {
	r.POST("", h.Create(i))
	r.PUT("/:id", h.Update(i))
	r.DELETE("/:id", h.Delete(i))
	r.GET("/:id", h.Get(i))
	r.GET("/app/:appId", h.ListByApp(i))
	r.GET("/:id/health", h.CheckHealth(i))
}

// Create godoc
// @Summary Create a new gateway
// @Description Create a new gateway
// @Tags Gateway
// @Accept json
// @Produce json
// @Param request body model.Gateway true "Gateway details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gateways [post]
func (h *Gateway) Create(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		var gateway model.Gateway
		if err := c.Bind(&gateway); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
		}

		if err := h.gatewayService.Create(c.Request().Context(), gateway); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
		}

		return c.Redirect(http.StatusSeeOther, "/dashboard/gateways")
	}
}

// Update godoc
// @Summary Update a gateway
// @Description Update a gateway
// @Tags Gateway
// @Accept json
// @Produce json
// @Param id path string true "Gateway ID"
// @Param request body model.Gateway true "Gateway details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gateways/{id} [put]
func (h *Gateway) Update(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var gateway model.Gateway
		if err := c.Bind(&gateway); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
		}

		gateway.ID = id
		if err := h.gatewayService.Update(c.Request().Context(), gateway); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
		}

		return c.Redirect(http.StatusSeeOther, "/dashboard/gateways")
	}
}

// Delete godoc
// @Summary Delete a gateway
// @Description Delete a gateway
// @Tags Gateway
// @Accept json
// @Produce json
// @Param id path string true "Gateway ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gateways/{id} [delete]
func (h *Gateway) Delete(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if err := h.gatewayService.Delete(c.Request().Context(), id); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
		}

		return c.Redirect(http.StatusSeeOther, "/dashboard/gateways")
	}
}

// Get godoc
// @Summary Get a gateway by ID
// @Description Get a gateway by ID
// @Tags Gateway
// @Accept json
// @Produce json
// @Param id path string true "Gateway ID"
// @Success 200 {object} model.Gateway
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gateways/{id} [get]
func (h *Gateway) Get(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		gateway, err := h.gatewayService.Get(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
		}

		return c.JSON(http.StatusOK, gateway)
	}
}

// ListByApp godoc
// @Summary List gateways by app ID
// @Description List gateways by app ID
// @Tags Gateway
// @Accept json
// @Produce json
// @Param appId path string true "App ID"
// @Success 200 {object} []model.Gateway
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gateways/app/{appId} [get]
func (h *Gateway) ListByApp(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		appID := c.Param("appId")
		gateways, err := h.gatewayService.ListByApp(c.Request().Context(), appID)
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
		}

		return c.JSON(http.StatusOK, gateways)
	}
}

// CheckHealth godoc
// @Summary Check the health of a gateway
// @Description Check the health of a gateway
// @Tags Gateway
// @Accept json
// @Produce json
// @Param id path string true "Gateway ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gateways/{id}/health [get]
func (h *Gateway) CheckHealth(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		gateway, err := h.gatewayService.Get(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
		}

		if err := h.healthChecker.CheckGatewayHealth(c.Request().Context(), gateway); err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, map[string]string{
				"status":  "unhealthy",
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"message": "Gateway is healthy",
		})
	}
}
