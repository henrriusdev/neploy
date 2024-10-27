package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	inertia "github.com/romsar/gonertia"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Auth struct {
	user service.User
}

func NewAuth(user service.User) *Auth {
	return &Auth{user}
}

func (a *Auth) RegisterRoutes(r fiber.Router, i *inertia.Inertia) {
	r.Post("/login", a.Login)
	r.Get("/logout", adaptor.HTTPHandler(a.Logout(i)))
	r.Get("", adaptor.HTTPHandler(a.Index(i)))
	r.Get("/onboard", adaptor.HTTPHandler(a.Onboard(i)))
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
