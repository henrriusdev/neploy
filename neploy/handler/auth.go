package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	inertia "github.com/romsar/gonertia"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/gitlab"
	"neploy.dev/config"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Auth struct {
	user service.User
	i    *inertia.Inertia
}

func NewAuth(user service.User, i *inertia.Inertia) *Auth {
	return &Auth{
		user: user,
		i:    i,
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

func (a *Auth) RegisterRoutes(r *echo.Group) {
	r.POST("/login", a.Login)
	r.GET("/logout", a.Logout)
	r.GET("", a.Index)
	r.GET("/onboard", a.Onboard)
	r.GET("/auth/github", a.GithubOAuth)
	r.GET("/auth/github/callback", a.GithubOAuthCallback)
	r.GET("/auth/gitlab", a.GitlabOAuth)
	r.GET("/auth/gitlab/callback", a.GitlabOAuthCallback)
}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags Auth, User
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login Request"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /login [post]
func (a *Auth) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request",
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate and authenticate user
	res, err := a.user.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Invalid credentials",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = res.Token
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary Logout a user
// @Description Logout a user
// @Tags Auth, User
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /logout [post]
func (a *Auth) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.MaxAge = -1
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)

	a.i.Redirect(c.Response(), c.Request(), "/")

	return nil
}

func (a *Auth) Index(c echo.Context) error {
	return a.i.Render(c.Response(), c.Request(), "Home/Login", inertia.Props{})
}

func (a *Auth) Onboard(c echo.Context) error {
	return a.i.Render(c.Response(), c.Request(), "Home/Onboard", inertia.Props{})
}

// GithubOAuth godoc
// @Summary Start GitHub OAuth flow
// @Description Start GitHub OAuth flow
// @Tags Auth
// @Produce json
// @Param state query string false "State parameter"
// @Success 302 {string} string "Redirects to GitHub OAuth flow"
// @Router /auth/github [get]
func (a *Auth) GithubOAuth(c echo.Context) error {
	githubConfig := GetConfig(model.Github)
	state := c.QueryParam("state") // Get state parameter (invitation token)
	logger.Info("Starting GitHub OAuth with state: %s", state)
	url := githubConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GithubOAuthCallback godoc
// @Summary Handle GitHub OAuth callback
// @Description Handle GitHub OAuth callback
// @Tags Auth
// @Produce json
// @Param state query string false "State parameter"
// @Param code query string false "Authorization code"
// @Success 200 {object} model.OAuthResponse
// @Failure 500 {object} map[string]interface{}
// @Router /auth/github/callback [get]
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

// GitlabOAuth godoc
// @Summary Start GitLab OAuth flow
// @Description Start GitLab OAuth flow
// @Tags Auth
// @Produce json
// @Param state query string false "State parameter"
// @Success 302 {string} string "Redirects to GitLab OAuth flow"
// @Router /auth/gitlab [get]
func (a *Auth) GitlabOAuth(c echo.Context) error {
	gitlabConfig := GetConfig(model.Gitlab)
	state := c.QueryParam("state") // Get state parameter (invitation token)
	logger.Info("Starting GitLab OAuth with state: %s", state)
	url := gitlabConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GitlabOAuthCallback godoc
// @Summary Handle GitLab OAuth callback
// @Description Handle GitLab OAuth callback
// @Tags Auth
// @Produce json
// @Param state query string false "State parameter"
// @Param code query string false "Authorization code"
// @Success 200 {object} model.OAuthResponse
// @Failure 500 {object} map[string]interface{}
// @Router /auth/gitlab/callback [get]
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
