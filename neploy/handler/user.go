package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type User struct {
	user service.User
}

func NewUser(user service.User) *User {
	return &User{user: user}
}

func (u *User) RegisterRoutes(app *fiber.App, i *gonertia.Inertia) {
	app.Post("/invite", u.InviteUser)
}

func (h *User) InviteUser(c *fiber.Ctx) error {
	var req model.InviteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Get team ID from authenticated user
	teamID := c.Locals("team_id").(string)
	req.TeamID = teamID

	// Validate request
	if req.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email is required")
	}
	if req.Role == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Role is required")
	}

	// Send invitation
	if err := h.user.InviteUser(c.Context(), req); err != nil {
		if err.Error() == "user already exists" {
			return fiber.NewError(fiber.StatusConflict, "User already exists in the system")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send invitation")
	}

	return c.JSON(fiber.Map{
		"message": "Invitation sent successfully",
	})
}
