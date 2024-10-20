package neploy

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Start() {
	i := initInertia()
	fmt.Println(i == nil)

	app := fiber.New(fiber.Config{
		Concurrency: 10,
	})

	app.Use(adaptor.HTTPMiddleware(i.Middleware))

	loginRoutes(app, i)

	app.Get("/build/assets/:filename", func(c *fiber.Ctx) error {
		filename := c.Params("filename")

		if strings.HasSuffix(filename, ".js") {
			c.Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(filename, ".css") {
			c.Set("Content-Type", "text/css")
		}

		return c.SendFile("./public/build/assets/" + filename)
	})

	app.Listen(":3000")
}
