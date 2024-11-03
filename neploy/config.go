package neploy

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
	"neploy.dev/neploy/middleware"
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

	repos := NewRepositories(npy)
	npy.Repositories = repos

	services := NewServices(npy)
	npy.Services = services

	app.Use(adaptor.HTTPMiddleware(i.Middleware))
	app.Use(middleware.OnboardingMiddleware(services.Onboard))

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
	user := service.NewUser(npy.Repositories)
	role := service.NewRole(npy.Repositories.Role, npy.Repositories.UserRole)
	metadata := service.NewMetadata(npy.Repositories.Metadata)
	onboard := service.NewOnboard(user, role, metadata)

	return service.Services{
		User:    user,
		Role:    role,
		Onboard: onboard,
	}
}

func NewRepositories(npy Neploy) repository.Repositories {
	metadata := repository.NewMetadata(npy.DB)
	role := repository.NewRole(npy.DB)
	user := repository.NewUser(npy.DB)
	userOauth := repository.NewUserOauth(npy.DB)
	userRole := repository.NewUserRole(npy.DB)

	return repository.Repositories{
		Metadata:  metadata,
		Role:      role,
		User:      user,
		UserOauth: userOauth,
		UserRole:  userRole,
	}
}

func NewHandlers(npy Neploy, i *gonertia.Inertia, app *fiber.App) {
	loginRoutes(app, i, npy)
	onboardRoutes(app, i, npy)
}
