package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type User struct {
	user service.User
}

func NewUser(user service.User) *User {
	return &User{user: user}
}

func (u *User) RegisterRoutes(app fiber.Router, i *gonertia.Inertia) {
	app.Post("/invite", u.InviteUser)
	app.Get("/invite/:token", adaptor.HTTPHandlerFunc(u.AcceptInvite(i)))
	app.Post("/complete-invite", u.CompleteInvite)
}

func (h *User) InviteUser(c *fiber.Ctx) error {
	var req model.InviteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Check if user already exists
	_, err := h.user.GetByEmail(c.Context(), req.Email)
	if err == nil {
		return fiber.NewError(fiber.StatusConflict, "User already exists in the system")
	}

	// Validate request
	if req.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email is required")
	}
	if req.Role == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Role is required")
	}

	// Send invitation
	if err := h.user.InviteUser(c.Context(), req); err != nil {
		logger.Error("error inviting user: %v", err)
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
	var req model.CompleteInviteRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	// Get oauth_id from cookie
	oauthID := c.Cookies("oauth_id")
	if oauthID == "" {
		oauthID = "no_oauth_id"
	}

	// delete oauth_id cookie
	c.ClearCookie("oauth_id")
	req.OauthID = oauthID

	// Accept the invitation
	invitation, err := u.user.AcceptInvitation(c.Context(), req.Token)
	if err != nil {
		return err
	}

	// Create user
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

	if err := u.user.Create(c.Context(), userReq, oauthID); err != nil {
		return err
	}

	// Add user to role
	if err := u.user.AddUserRole(c.Context(), req.Email, invitation.Role); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (u *User) AcceptInvite(i *gonertia.Inertia) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// the token is a path variable
		token := strings.Split(r.URL.Path, "/invite/")[1]

		// Get OAuth callback data if present
		username := r.URL.Query().Get("username")
		email := r.URL.Query().Get("email")
		provider := r.URL.Query().Get("provider")

		// Get oauth_id from cookies
		var oauthID string
		for _, cookie := range r.Cookies() {
			if cookie.Name == "oauth_id" {
				oauthID = cookie.Value
				// Clear the cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "oauth_id",
					Value:    "",
					Path:     "/",
					Expires:  time.Now().Add(-1 * time.Hour),
					HttpOnly: true,
				})
				break
			}
		}

		// Obtener la invitación
		invitation, err := u.user.GetInvitationByToken(context.Background(), token)
		if err != nil {
			logger.Error("failed to get invitation: token=%s, error=%v", token, err)
			i.Render(w, r, "Auth/CompleteInvite", gonertia.Props{
				"token":  token,
				"error":  "Invalid or expired invitation",
				"status": "invalid",
			})
			return
		}

		props := gonertia.Props{
			"token":  token,
			"email":  invitation.Email,
			"status": "valid",
		}

		// Add OAuth data if present
		if username != "" {
			props["username"] = username
		}
		if email != "" {
			props["email"] = email
		}
		if provider != "" {
			props["provider"] = provider
		}
		if oauthID != "" {
			props["oauth_id"] = oauthID
		}

		i.Render(w, r, "Auth/CompleteInvite", props)
	}

	return fn
}
