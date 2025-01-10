package neploy

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	neployware "neploy.dev/neploy/middleware"
	neployway "neploy.dev/pkg/gateway"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/repository"
	"neploy.dev/pkg/service"
	"neploy.dev/pkg/store"
	"neploy.dev/pkg/websocket"

	_ "neploy.dev/neploy/docs"
)

type Neploy struct {
	DB           store.Queryable
	Port         string
	Services     service.Services
	Repositories repository.Repositories
	Router       *neployway.Router
}

func Start(npy Neploy) {
	i := initInertia()

	e := echo.New()

	// Initialize repositories
	repos := NewRepositories(npy)
	npy.Repositories = repos

	// Initialize router
	router := neployway.NewRouter(npy.Repositories.ApplicationStat)
	npy.Router = router

	// Initialize services
	services := NewServices(npy)
	npy.Services = services

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${remote_ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))
	e.Use(echo.WrapMiddleware(i.Middleware))

	// WebSocket routes with specialized handlers
	e.GET("/ws/notifications", websocket.UpgradeProgressWS())
	e.GET("/ws/interactive", websocket.UpgradeInteractiveWS())

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(neployware.OnboardingMiddleware(services.Onboard))
	logger.SetLogger()

	// Validator
	vldtr := validator.New()
	e.Validator = &CustomValidator{validator: vldtr}

	// Routes
	RegisterRoutes(e, i, npy)

	// Static files
	e.GET("/build/assets/:filename", func(c echo.Context) error {
		filename := c.Param("filename")

		if strings.HasSuffix(filename, ".js") {
			c.Response().Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(filename, ".css") {
			c.Response().Header().Set("Content-Type", "text/css")
		}

		return c.File("./public/build/assets/" + filename)
	})

	e.Start(":" + npy.Port)
}

func NewServices(npy Neploy) service.Services {
	application := service.NewApplication(npy.Repositories, npy.Router)
	metadata := service.NewMetadata(npy.Repositories.Metadata)
	user := service.NewUser(npy.Repositories)
	role := service.NewRole(npy.Repositories.Role, npy.Repositories.UserRole)
	onboard := service.NewOnboard(user, role, metadata)
	gateway := service.NewGateway(npy.Repositories.Gateway, npy.Repositories.Application, npy.Repositories.ApplicationStat)

	return service.Services{
		Application: application,
		Gateway:     gateway,
		Metadata:    metadata,
		Onboard:     onboard,
		Role:        role,
		User:        user,
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
	techStack := repository.NewTechStack(npy.DB)
	gateway := repository.NewGateway(npy.DB)

	return repository.Repositories{
		Application:     application,
		ApplicationStat: applicationStat,
		Gateway:         gateway,
		Metadata:        metadata,
		Role:            role,
		TechStack:       techStack,
		User:            user,
		UserOauth:       userOauth,
		UserRole:        userRole,
		UserTechStack:   userTechStack,
		VisitorInfo:     visitorInfo,
		VisitorTrace:    visitorTrace,
	}
}
