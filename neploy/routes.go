package neploy

import (
	"github.com/gofiber/fiber/v2"
	"neploy.dev/neploy/handler"

	inertia "github.com/romsar/gonertia"
)

func loginRoutes(app *fiber.App, i *inertia.Inertia) {
	auth := handler.NewAuth()
	auth.RegisterRoutes(app.Group(""), i)
}
