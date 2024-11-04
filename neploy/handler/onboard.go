package handler

import (
	"strconv"

	"neploy.dev/neploy/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Onboard struct {
	validator validation.XValidator
	service   service.Onboard
}

func NewOnboard(validator validation.XValidator, service service.Onboard) *Onboard {
	return &Onboard{
		validator: validator,
		service:   service,
	}
}

func (o *Onboard) RegisterRoutes(r fiber.Router, i *gonertia.Inertia) {
	r.Post("", o.Initiate)
}

func (o *Onboard) Initiate(c *fiber.Ctx) error {
	var req model.OnboardRequest
	if err := c.BodyParser(&req); err != nil {
		log.Err(err).Msg("error")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	oauthID := c.Cookies("oauth_id")
	if oauthID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	oauth, _ := strconv.Atoi(oauthID)
	req.OauthID = oauth

	if err := o.service.Initiate(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Onboarding initiated"})
}
