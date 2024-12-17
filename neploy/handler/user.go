package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
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
	app.Get("/invite/:token", adaptor.HTTPHandlerFunc(u.AcceptInvite(i)))
	app.Post("/users/complete-invite", u.CompleteInvite)
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

	return c.JSON(gonertia.Props{
		"message": "Invitation sent successfully",
	})
}

func (u *User) CompleteInvite(c *fiber.Ctx) error {
	var req struct {
		Token     string    `json:"token"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		DOB       time.Time `json:"dob"`
		Phone     string    `json:"phone"`
		Address   string    `json:"address"`
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		Password  string    `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	// Aceptar la invitación
	if err := u.user.AcceptInvitation(c.Context(), req.Token); err != nil {
		return err
	}

	// Crear el usuario
	userReq := model.CreateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DOB:       model.Date{Time: req.DOB},
		Phone:     req.Phone,
		Address:   req.Address,
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
	}

	if err := u.user.Create(c.Context(), userReq, 0); err != nil {
		return err
	}

	return c.JSON(gonertia.Props{
		"message": "User created successfully",
	})
}

func (u *User) AcceptInvite(i *gonertia.Inertia) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		// Obtener la invitación
		invitation, err := u.user.GetInvitationByToken(context.Background(), token)
		if err != nil {
			i.Render(w, r, "Auth/AcceptInvite", gonertia.Props{
				"token":   token,
				"expired": true,
			})
		}

		// Verificar estado
		if time.Now().After(invitation.ExpiresAt.Time) {
			i.Render(w, r, "Auth/AcceptInvite", gonertia.Props{
				"token":   token,
				"expired": true,
			})
		}

		if invitation.AcceptedAt != nil {
			i.Render(w, r, "Auth/AcceptInvite", gonertia.Props{
				"token":           token,
				"alreadyAccepted": true,
			})
		}

		// Redirigir al flujo de completar invitación
		i.Render(w, r, "Auth/CompleteInvite", gonertia.Props{
			"token": token,
			"email": invitation.Email,
		})
	}

	return fn
}
