package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	inertia "github.com/romsar/gonertia"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/gitlab"
	"neploy.dev/config"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Auth struct {
	user service.User
}

func NewAuth(user service.User) *Auth {
	return &Auth{user}
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

func (a *Auth) RegisterRoutes(r fiber.Router, i *inertia.Inertia) {
	r.Post("/login", a.Login)
	r.Get("/logout", adaptor.HTTPHandler(a.Logout(i)))
	r.Get("", adaptor.HTTPHandler(a.Index(i)))
	r.Get("/onboard", adaptor.HTTPHandler(a.Onboard(i)))
	r.Get("/auth/github", a.GithubOAuth)
	r.Get("/auth/github/callback", a.GithubOAuthCallback)
	r.Get("/auth/gitlab", a.GitlabOAuth)
	r.Get("/auth/gitlab/callback", a.GitlabOAuthCallback)
}

func (a *Auth) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// todo: validate the request

	res, err := a.user.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// put the token in the session
	// c.Session.Set("token", res.Token)

	return c.JSON(fiber.Map{"token": res.Token})
}

func (a *Auth) Logout(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// logout logic here
	}
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
	url := githubConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (a *Auth) GithubOAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	githubConfig := GetConfig(model.Github)
	token, err := githubConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	client := githubConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
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
		HTTPOnly: true,
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

	// Si el correo aún es vacío, maneja el caso de no poder obtener un correo electrónico
	if user.Email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("No email available for this user")
	}

	oauthResponse := model.OAuthResponse{
		Provider: model.Github,
		Email:    user.Email,
		Username: user.Login,
	}

	return c.Redirect(fmt.Sprintf("/onboard?username=%s&email=%s&provider=%s", oauthResponse.Username, oauthResponse.Email, oauthResponse.Provider))
}

func (a *Auth) GitlabOAuth(c *fiber.Ctx) error {
	gitlabConfig := GetConfig(model.Gitlab)
	url := gitlabConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (a *Auth) GitlabOAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	gitlabConfig := GetConfig(model.Gitlab)
	token, err := gitlabConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	client := gitlabConfig.Client(context.Background(), token)
	resp, err := client.Get("https://gitlab.com/api/v4/user")
	if err != nil {
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

	oauthResponse := model.OAuthResponse{
		Provider: model.Gitlab,
		Email:    user.Email,
		Username: user.Username,
	}

	return c.Redirect(fmt.Sprintf("/onboard?username=%s&email=%s&provider=%s", oauthResponse.Username, oauthResponse.Email, oauthResponse.Provider))
}
