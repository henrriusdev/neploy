package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
}

func NewAuth(validator validation.XValidator, user service.User) *Auth {
	return &Auth{
		validator: validator,
		user:      user,
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
	r.POST("/login", a.Login(i))
	r.GET("/logout", a.Logout(i))
	r.GET("", a.Index(i))
	r.GET("/onboard", a.Onboard(i))
	r.GET("/auth/github", a.GithubOAuth)
	r.GET("/auth/github/callback", a.GithubOAuthCallback)
	r.GET("/auth/gitlab", a.GitlabOAuth)
	r.GET("/auth/gitlab/callback", a.GitlabOAuthCallback)
}

func (a *Auth) Login(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
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

			return echo.NewHTTPError(http.StatusBadRequest, strings.Join(errMsgs, " and "))
		}

		// Validate and authenticate user
		res, err := a.user.Login(c.Request().Context(), req)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid credentials",
			})
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.JWTClaims{
			ID:       res.User.ID,
			Email:    res.User.Email,
			Name:     res.User.FirstName + " " + res.User.LastName,
			Username: res.User.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		})

		tokenString, err := token.SignedString([]byte(config.Env.JWTSecret))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to generate token",
			})
		}

		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = tokenString
		cookie.HttpOnly = true
		cookie.Path = "/"
		c.SetCookie(cookie)

		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}
}

func (h *Auth) Logout(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = ""
		cookie.HttpOnly = true
		cookie.Path = "/"
		cookie.MaxAge = -1
		cookie.Expires = time.Unix(0, 0)
		c.SetCookie(cookie)

		i.Redirect(c.Response(), c.Request(), "/")

		return nil
	}
}

func (a *Auth) Index(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		return i.Render(c.Response(), c.Request(), "Home/Login", inertia.Props{})
	}
}

func (a *Auth) Onboard(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		return i.Render(c.Response(), c.Request(), "Home/Onboard", inertia.Props{})
	}
}

func (a *Auth) GithubOAuth(c echo.Context) error {
	githubConfig := GetConfig(model.Github)
	state := c.QueryParam("state") // Get state parameter (invitation token)
	logger.Info("Starting GitHub OAuth with state: %s", state)
	url := githubConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *Auth) GithubOAuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	logger.Info("GitHub OAuth callback received with state: %s", state)
	githubConfig := GetConfig(model.Github)
	token, err := githubConfig.Exchange(context.Background(), code)
	if err != nil {
		logger.Error("Failed to exchange token: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to exchange token")
	}

	client := githubConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		logger.Error("Failed to get user info: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to get user info")
	}
	defer resp.Body.Close()

	var user struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to parse user info")
	}

	// Set oauth_id cookie
	cookie := new(http.Cookie)
	cookie.Name = "oauth_id"
	cookie.Value = fmt.Sprintf("%d", user.ID)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteLaxMode
	cookie.MaxAge = 24 * 60 * 60
	cookie.Domain = c.Request().Host
	c.SetCookie(cookie)

	if user.Email == "" {
		resp, err = client.Get("https://api.github.com/user/emails")
		if err != nil {
			logger.Error("Failed to get user emails: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to get user emails")
		}
		defer resp.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to parse user emails")
		}

		// Find primary and verified email
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
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/users/invite/%s?username=%s&email=%s&provider=%s",
			state,
			oauthResponse.Username,
			oauthResponse.Email,
			oauthResponse.Provider))
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/onboard?username=%s&email=%s&provider=%s",
		oauthResponse.Username,
		oauthResponse.Email,
		oauthResponse.Provider))
}

func (a *Auth) GitlabOAuth(c echo.Context) error {
	gitlabConfig := GetConfig(model.Gitlab)
	state := c.QueryParam("state") // Get state parameter (invitation token)
	logger.Info("Starting GitLab OAuth with state: %s", state)
	url := gitlabConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *Auth) GitlabOAuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	logger.Info("GitLab OAuth callback received with state: %s", state)
	gitlabConfig := GetConfig(model.Gitlab)
	token, err := gitlabConfig.Exchange(context.Background(), code)
	if err != nil {
		logger.Error("Failed to exchange token: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to exchange token")
	}

	client := gitlabConfig.Client(context.Background(), token)
	resp, err := client.Get("https://gitlab.com/api/v4/user")
	if err != nil {
		logger.Error("Failed to get user info: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to get user info")
	}
	defer resp.Body.Close()

	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to parse user info")
	}

	// Set oauth_id cookie
	cookie := new(http.Cookie)
	cookie.Name = "oauth_id"
	cookie.Value = token.AccessToken
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteLaxMode
	cookie.MaxAge = 24 * 60 * 60
	cookie.Domain = c.Request().Host
	c.SetCookie(cookie)

	oauthResponse := model.OAuthResponse{
		Username: user.Username,
		Email:    user.Email,
		Provider: model.Gitlab,
	}

	if state != "" {
		// If we have an invitation token, redirect to the invite completion
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/users/invite/%s?username=%s&email=%s&provider=%s",
			state,
			oauthResponse.Username,
			oauthResponse.Email,
			oauthResponse.Provider))
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/onboard?username=%s&email=%s&provider=%s",
		oauthResponse.Username,
		oauthResponse.Email,
		oauthResponse.Provider))
}
