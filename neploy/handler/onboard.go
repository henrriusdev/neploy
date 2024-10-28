package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Onboard struct {
	service service.Onboard
}

func NewOnboard(service service.Onboard) *Onboard {
	return &Onboard{service}
}

func (o *Onboard) RegisterRoutes(r fiber.Router, i *gonertia.Inertia) {
	r.Post("/onboard", o.Initiate)
}

func (o *Onboard) Initiate(c *fiber.Ctx) error {
	var req model.OnboardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := o.service.Initiate(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Onboarding initiated"})
}
