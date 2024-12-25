package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"

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

func (h *Gateway) RegisterRoutes(r fiber.Router, i *gonertia.Inertia) {
	r.Post("", h.Create)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
	r.Get("/:id", h.Get)
	r.Get("/app/:appId", h.ListByApp)
}

func (h *Gateway) Create(c *fiber.Ctx) error {
	var gateway model.Gateway
	if err := c.BodyParser(&gateway); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.gatewayService.Create(c.Context(), gateway); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(gateway)
}

func (h *Gateway) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var gateway model.Gateway
	if err := c.BodyParser(&gateway); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	gateway.ID = id
	if err := h.gatewayService.Update(c.Context(), gateway); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(gateway)
}

func (h *Gateway) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.gatewayService.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Gateway) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	gateway, err := h.gatewayService.Get(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(gateway)
}

func (h *Gateway) ListByApp(c *fiber.Ctx) error {
	appID := c.Params("appId")
	gateways, err := h.gatewayService.ListByApp(c.Context(), appID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(gateways)
}
