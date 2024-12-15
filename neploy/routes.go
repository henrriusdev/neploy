package neploy

import (
	"github.com/gofiber/fiber/v2"
	"neploy.dev/neploy/handler"

	inertia "github.com/romsar/gonertia"
)

func loginRoutes(app *fiber.App, i *inertia.Inertia, npy Neploy) {
	auth := handler.NewAuth(npy.Validator, npy.Services.User, npy.SessionStore)
	auth.RegisterRoutes(app.Group(""), i)
}

func onboardRoutes(app *fiber.App, i *inertia.Inertia, npy Neploy) {
	onboard := handler.NewOnboard(npy.Validator, npy.Services.Onboard)
	onboard.RegisterRoutes(app.Group("/onboard"), i)
}

func dashboardRoutes(app *fiber.App, i *inertia.Inertia, npy Neploy) {
	dashboard := handler.NewDashboard(npy.Services.Metadata, npy.Services.Application, npy.Services.User, npy.SessionStore)
	dashboard.RegisterRoutes(app.Group("/dashboard"), i)
}

func userRoutes(app *fiber.App, npy Neploy) {
	h := handler.NewUser(npy.Services.User)
	users := app.Group("/users")
	users.Post("/invite", h.InviteUser)
}

func RegisterRoutes(app *fiber.App, i *inertia.Inertia, npy Neploy) {
	loginRoutes(app, i, npy)
	onboardRoutes(app, i, npy)
	dashboardRoutes(app, i, npy)
	userRoutes(app, npy)
}
