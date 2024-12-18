package neploy

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"neploy.dev/neploy/middleware"
	"neploy.dev/neploy/validation"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/repository"
	"neploy.dev/pkg/service"
	"neploy.dev/pkg/store"
)

type Neploy struct {
	DB           store.Queryable
	Port         string
	Services     service.Services
	Repositories repository.Repositories
	Validator    validation.XValidator
	SessionStore *session.Store // Add session store to Neploy struct
}

// Create a new config struct for session settings
type SessionConfig struct {
	Expiration     time.Duration
	CookieName     string
	CookieSecure   bool
	CookieHTTPOnly bool
}

// Initialize session store with default config
func NewSessionStore() *session.Store {
	return session.New(session.Config{
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:session",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})
}

func Start(npy Neploy) {
	// Initialize session store
	sessionStore := NewSessionStore()
	npy.SessionStore = sessionStore

	i := initInertia()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(validation.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
		Concurrency: 10,
	})

	// Set the global logger

	// Initialize repositories
	repos := NewRepositories(npy)
	npy.Repositories = repos

	// Initialize services
	services := NewServices(npy)
	npy.Services = services

	// Middleware
	app.Use(flogger.New(flogger.Config{
		Format:     "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(adaptor.HTTPMiddleware(i.Middleware))
	app.Use(middleware.OnboardingMiddleware(services.Onboard))
	app.Use(middleware.SessionMiddleware(npy.SessionStore))
	logger.SetLogger()

	// Validator
	myValidator := &validation.XValidator{
		Validator: validation.Validate,
	}
	npy.Validator = *myValidator

	// Routes
	RegisterRoutes(app, i, npy)

	// Static files
	app.Get("/build/assets/:filename", func(c *fiber.Ctx) error {
		filename := c.Params("filename")

		if strings.HasSuffix(filename, ".js") {
			c.Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(filename, ".css") {
			c.Set("Content-Type", "text/css")
		}

		return c.SendFile("./public/build/assets/" + filename)
	})

	// print all routes
	app.Get("/routes", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"routes": app.GetRoutes()})
	})
	app.Listen(":" + npy.Port)
}

func NewServices(npy Neploy) service.Services {
	application := service.NewApplication(npy.Repositories.Application, npy.Repositories.ApplicationStat)
	metadata := service.NewMetadata(npy.Repositories.Metadata)
	email := service.NewEmail()
	user := service.NewUser(npy.Repositories, email)
	role := service.NewRole(npy.Repositories.Role, npy.Repositories.UserRole)
	onboard := service.NewOnboard(user, role, metadata)

	return service.Services{
		Application: application,
		User:        user,
		Role:        role,
		Metadata:    metadata,
		Onboard:     onboard,
	}
}

func NewRepositories(npy Neploy) repository.Repositories {
	metadata := repository.NewMetadata(npy.DB)
	role := repository.NewRole(npy.DB)
	user := repository.NewUser(npy.DB)
	userOauth := repository.NewUserOauth(npy.DB)
	userRole := repository.NewUserRole(npy.DB)
	application := repository.NewApplication(npy.DB)
	applicationStat := repository.NewApplicationStat(npy.DB)
	userTechStack := repository.NewUserTechStack(npy.DB)
	visitorInfo := repository.NewVisitor(npy.DB)
	visitorTrace := repository.NewVisitorTrace(npy.DB)

	return repository.Repositories{
		Metadata:        metadata,
		Role:            role,
		User:            user,
		UserOauth:       userOauth,
		UserRole:        userRole,
		Application:     application,
		ApplicationStat: applicationStat,
		UserTechStack:   userTechStack,
		VisitorInfo:     visitorInfo,
		VisitorTrace:    visitorTrace,
	}
}
