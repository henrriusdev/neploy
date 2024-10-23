package neploy

import (
	"github.com/gofiber/fiber/v2"
	"neploy.dev/neploy/handler"

	inertia "github.com/romsar/gonertia"
)

func loginRoutes(app *fiber.App, i *inertia.Inertia, npy Neploy) {
	auth := handler.NewAuth(npy.Services.User)
	auth.RegisterRoutes(app.Group(""), i)
}
