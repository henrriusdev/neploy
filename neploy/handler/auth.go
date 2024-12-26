package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labstack/echo/v4"
	inertia "github.com/romsar/gonertia"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/gitlab"
	"neploy.dev/config"
	"neploy.dev/neploy/validation"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Auth struct {
	validator validation.XValidator
	user      service.User
	sessions  *session.Store
}

func NewAuth(validator validation.XValidator, user service.User, session *session.Store) *Auth {
	return &Auth{
		validator: validator,
		user:      user,
		sessions:  session,
	}
}

func GetConfig(provider model.Provider) *oauth2.Config {
	switch provider {
	case model.Github:
		return &oauth2.Config{
			ClientID:     config.Env.GithubClientID,
			ClientSecret: config.Env.GithubClientSecret,
			RedirectURL:  fmt.Sprintf("%s:%s/auth/github/callback", config.Env.BaseURL, config.Env.Port),
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		}
	case model.Gitlab:
		return &oauth2.Config{
			ClientID:     config.Env.GitlabApplicationID,
			ClientSecret: config.Env.GitlabSecret,
			RedirectURL:  fmt.Sprintf("%s:%s/auth/gitlab/callback", config.Env.BaseURL, config.Env.Port),
			Scopes:       []string{"read_user"},
			Endpoint:     gitlab.Endpoint,
		}
	default:
		return nil
	}
}

func (a *Auth) RegisterRoutes(r *echo.Group, i *inertia.Inertia) {
	r.POST("/login", a.Login)
	r.GET("/logout", a.Logout)
	r.GET("", adaptor.HTTPHandler(a.Index(i)))
	r.GET("/onboard", adaptor.HTTPHandler(a.Onboard(i)))
	r.GET("/auth/github", a.GithubOAuth)
	r.GET("/auth/github/callback", a.GithubOAuthCallback)
	r.GET("/auth/gitlab", a.GitlabOAuth)
	r.GET("/auth/gitlab/callback", a.GitlabOAuthCallback)
}

func (a *Auth) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if errs := a.validator.Validate(req); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}

	// Validate and authenticate user
	res, err := a.user.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Get session from store
	sess, err := a.sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}

	// Set session values
	sess.Set("authenticated", true)
	sess.Set("user_id", res.User.ID)
	sess.Set("username", res.User.Username)
	sess.Set("email", res.User.Email)
	sess.Set("name", res.User.FirstName+" "+res.User.LastName)
	sess.Set("token", res.Token)

	// put a cookie with token and user_id
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    sess.ID(),
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    res.Token,
		HTTPOnly: true,
	})

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save session",
		})
	}

	return c.JSON(fiber.Map{
		"token": res.Token,
	})
}

func (h *Auth) Logout(c *fiber.Ctx) error {
	// Get the session
	sess, err := h.sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}

	// Clear all session data
	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout",
		})
	}

	// Optional: Clear the cookie explicitly
	c.ClearCookie("session")

	// Redirect to login page or return success response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged out",
	})
}

func (a *Auth) Index(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i.Render(w, r, "Home/Login", inertia.Props{})
	}
}

func (a *Auth) Onboard(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i.Render(w, r, "Home/Onboard", inertia.Props{})
	}
}

func (a *Auth) GithubOAuth(c *fiber.Ctx) error {
	githubConfig := GetConfig(model.Github)
	state := c.Query("state") // Get state parameter (invitation token)
	logger.Info("Starting GitHub OAuth with state: %s", state)
	url := githubConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (a *Auth) GithubOAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	logger.Info("GitHub OAuth callback received with state: %s", state)
	githubConfig := GetConfig(model.Github)
	token, err := githubConfig.Exchange(context.Background(), code)
	if err != nil {
		logger.Error("Failed to exchange token: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	client := githubConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		logger.Error("Failed to get user info: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer resp.Body.Close()

	var user struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user info")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_id",
		Value:    fmt.Sprintf("%d", user.ID),
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   24 * 60 * 60,
		Domain:   c.Hostname(),
	})

	if user.Email == "" {
		resp, err = client.Get("https://api.github.com/user/emails")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user emails")
		}
		defer resp.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user emails")
		}

		// Encuentra el correo principal y verificado
		for _, e := range emails {
			if e.Primary && e.Verified {
				user.Email = e.Email
				break
			}
		}
	}

	oauthResponse := model.OAuthResponse{
		Provider: model.Github,
		Email:    user.Email,
		Username: user.Login,
	}

	if state != "" {
		// If we have an invitation token, redirect to the invite completion
		return c.Redirect(fmt.Sprintf("/users/invite/%s?username=%s&email=%s&provider=%s",
			state,
			oauthResponse.Username,
			oauthResponse.Email,
			oauthResponse.Provider))
	}

	return c.Redirect(fmt.Sprintf("/onboard?username=%s&email=%s&provider=%s",
		oauthResponse.Username,
		oauthResponse.Email,
		oauthResponse.Provider))
}

func (a *Auth) GitlabOAuth(c *fiber.Ctx) error {
	gitlabConfig := GetConfig(model.Gitlab)
	state := c.Query("state") // Get state parameter (invitation token)
	logger.Info("Starting GitLab OAuth with state: %s", state)
	url := gitlabConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (a *Auth) GitlabOAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	logger.Info("GitLab OAuth callback received with state: %s", state)
	gitlabConfig := GetConfig(model.Gitlab)
	token, err := gitlabConfig.Exchange(context.Background(), code)
	if err != nil {
		logger.Error("Failed to exchange token: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	client := gitlabConfig.Client(context.Background(), token)
	resp, err := client.Get("https://gitlab.com/api/v4/user")
	if err != nil {
		logger.Error("Failed to get user info: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer resp.Body.Close()

	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user info")
	}

	// Set oauth_id cookie with the access token
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_id",
		Value:    token.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   24 * 60 * 60,
		Domain:   c.Hostname(),
	})

	oauthResponse := model.OAuthResponse{
		Username: user.Username,
		Email:    user.Email,
		Provider: model.Gitlab,
	}

	if state != "" {
		// If we have an invitation token, redirect to the invite completion
		return c.Redirect(fmt.Sprintf("/users/invite/%s?username=%s&email=%s&provider=%s",
			state,
			oauthResponse.Username,
			oauthResponse.Email,
			oauthResponse.Provider))
	}

	return c.Redirect(fmt.Sprintf("/onboard?username=%s&email=%s&provider=%s",
		oauthResponse.Username,
		oauthResponse.Email,
		oauthResponse.Provider))
}
