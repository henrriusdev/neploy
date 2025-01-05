package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	inertia "github.com/romsar/gonertia"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type User struct {
	user service.User
	i    *inertia.Inertia
}

func NewUser(user service.User, i *inertia.Inertia) *User {
	return &User{user: user, i: i}
}

func (u *User) RegisterRoutes(r *echo.Group) {
	r.POST("/invite", u.InviteUser)
	r.GET("/invite/:token", u.AcceptInvite)
	r.POST("/complete-invite", u.CompleteInvite)
}

// InviteUser godoc
// @Summary Invite a user
// @Description Invite a user
// @Tags User
// @Accept json
// @Produce json
// @Param request body model.InviteUserRequest true "Invite User Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/invite [post]
func (h *User) InviteUser(c echo.Context) error {
	var req model.InviteUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Check if user already exists
	_, err := h.user.GetByEmail(c.Request().Context(), req.Email)
	if err == nil {
		return echo.NewHTTPError(http.StatusConflict, "User already exists in the system")
	}

	// Validate request
	if req.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is required")
	}
	if req.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Role is required")
	}

	// Send invitation
	if err := h.user.InviteUser(c.Request().Context(), req); err != nil {
		logger.Error("error inviting user: %v", err)
		if err.Error() == "user already exists" {
			return echo.NewHTTPError(http.StatusConflict, "User already exists in the system")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send invitation")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invitation sent successfully",
	})
}

// CompleteInvite godoc
// @Summary Complete user invitation
// @Description Complete user invitation
// @Tags User
// @Accept json
// @Produce json
// @Param request body model.CompleteInviteRequest true "Complete Invite Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/complete-invite [post]
func (u *User) CompleteInvite(c echo.Context) error {
	var req model.CompleteInviteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Get oauth_id from cookie
	cookie, err := c.Cookie("oauth_id")
	oauthID := "no_oauth_id"
	if err == nil && cookie.Value != "" {
		oauthID = cookie.Value
	}

	// delete oauth_id cookie
	cookieDel := new(http.Cookie)
	cookieDel.Name = "oauth_id"
	cookieDel.Value = ""
	cookieDel.Path = "/"
	cookieDel.MaxAge = -1
	c.SetCookie(cookieDel)

	req.OauthID = oauthID

	// Accept the invitation
	invitation, err := u.user.AcceptInvitation(c.Request().Context(), req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	if err := u.user.Create(c.Request().Context(), userReq, oauthID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Add user to role
	if err := u.user.AddUserRole(c.Request().Context(), req.Email, invitation.Role); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (u *User) AcceptInvite(c echo.Context) error {
	token := c.Param("token")
	username := c.QueryParam("username")
	email := c.QueryParam("email")
	provider := c.QueryParam("provider")

	// Get oauth_id from cookies
	cookie, err := c.Cookie("oauth_id")
	var oauthID string
	if err == nil && cookie.Value != "" {
		oauthID = cookie.Value
		// Clear the cookie
		cookieDel := new(http.Cookie)
		cookieDel.Name = "oauth_id"
		cookieDel.Value = ""
		cookieDel.Path = "/"
		cookieDel.Expires = time.Now().Add(-1 * time.Hour)
		cookieDel.HttpOnly = true
		c.SetCookie(cookieDel)
	}

	// Get the invitation
	invitation, err := u.user.GetInvitationByToken(context.Background(), token)
	if err != nil {
		logger.Error("failed to get invitation: token=%s, error=%v", token, err)
		return u.i.Render(c.Response(), c.Request(), "Auth/CompleteInvite", inertia.Props{
			"token":  token,
			"error":  "Invalid or expired invitation",
			"status": "invalid",
		})
	}

	props := inertia.Props{
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

	return u.i.Render(c.Response(), c.Request(), "Auth/CompleteInvite", props)
}
