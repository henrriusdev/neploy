package neploy

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
	"neploy.dev/pkg/repository"
	"neploy.dev/pkg/service"
	"neploy.dev/pkg/store"
)

type Neploy struct {
	DB           store.Queryable
	Port         string
	Services     service.Services
	Repositories repository.Repositories
}

func Start(npy Neploy) {
	i := initInertia()

	app := fiber.New(fiber.Config{
		Concurrency: 10,
	})

	app.Use(adaptor.HTTPMiddleware(i.Middleware))

	services := NewServices(npy)
	repos := NewRepositories(npy)

	npy.Services = services
	npy.Repositories = repos
	NewHandlers(npy, i, app)

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

func NewServices(npy Neploy) service.Services {
	user := service.NewUser(npy.Repositories.User, npy.Repositories.UserRole)
	return service.Services{
		User: user,
	}
}

func NewRepositories(npy Neploy) repository.Repositories {
	userRepo := repository.NewUser(npy.DB)
	return repository.Repositories{
		User: userRepo,
	}
}

func NewHandlers(npy Neploy, i *gonertia.Inertia, app *fiber.App) {
	loginRoutes(app, i, npy)
}
