package neploy

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	neployware "neploy.dev/neploy/middleware"
	"neploy.dev/neploy/validation"
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
	Validator    validation.XValidator
}

// Create a new config struct for session settings
type SessionConfig struct {
	Expiration     time.Duration
	CookieName     string
	CookieSecure   bool
	CookieHTTPOnly bool
}

func Start(npy Neploy) {
	i := initInertia()

	e := echo.New()

	// Custom error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			if err := c.JSON(he.Code, validation.GlobalErrorHandlerResp{
				Success: false,
				Message: he.Message.(string),
			}); err != nil {
				e.Logger.Error(err)
			}
		} else {
			if err := c.JSON(http.StatusBadRequest, validation.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			}); err != nil {
				e.Logger.Error(err)
			}
		}
	}

	// Initialize repositories
	repos := NewRepositories(npy)
	npy.Repositories = repos

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
	myValidator := &validation.XValidator{
		Validator: validation.Validate,
	}
	npy.Validator = *myValidator

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
	application := service.NewApplication(npy.Repositories.Application, npy.Repositories.ApplicationStat, npy.Repositories.TechStack, npy.Repositories.Gateway)
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
	techStack := repository.NewTechStack(npy.DB)

	return repository.Repositories{
		Application:     application,
		ApplicationStat: applicationStat,
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
