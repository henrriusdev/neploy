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
	dashboard := handler.NewDashboard(npy.Services.Metadata, npy.Services.Application, npy.SessionStore)
	dashboard.RegisterRoutes(app.Group("/dashboard"), i)
}
