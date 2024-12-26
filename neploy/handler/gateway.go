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
}

func NewGateway(gatewayService service.Gateway) *Gateway {
	return &Gateway{
		gatewayService: gatewayService,
	}
}

func (h *Gateway) RegisterRoutes(r *echo.Group, i *inertia.Inertia) {
	r.POST("", h.Create(i))
	r.PUT("/:id", h.Update(i))
	r.DELETE("/:id", h.Delete(i))
	r.GET("/:id", h.Get(i))
	r.GET("/app/:appId", h.ListByApp(i))
}

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

func (h *Gateway) Delete(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if err := h.gatewayService.Delete(c.Request().Context(), id); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
		}

		return c.Redirect(http.StatusSeeOther, "/dashboard/gateways")
	}
}

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
