package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	inertia "github.com/romsar/gonertia"
	"neploy.dev/pkg/service"
)

type Auth struct {
	user service.User
}

func NewAuth(user service.User) *Auth {
	return &Auth{user}
}

func (a *Auth) RegisterRoutes(r fiber.Router, i *inertia.Inertia) {
	r.Post("/login", adaptor.HTTPHandler(a.Login(i)))
	r.Get("/logout", adaptor.HTTPHandler(a.Logout(i)))
	r.Get("", adaptor.HTTPHandler(a.Index(i)))
}

func (a *Auth) Login(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// login logic here
	}
}

func (a *Auth) Logout(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// logout logic here
	}
}

func (a *Auth) Index(i *inertia.Inertia) http.HandlerFunc {
	fmt.Println("Index")
	return func(w http.ResponseWriter, r *http.Request) {
		i.Render(w, r, "Home/Login", inertia.Props{})
	}
}
