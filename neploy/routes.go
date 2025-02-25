package neploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"neploy.dev/neploy/handler"
	"neploy.dev/neploy/middleware"
	neployway "neploy.dev/pkg/gateway"
	"neploy.dev/pkg/logger"

	inertia "github.com/romsar/gonertia"
)

func loginRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	auth := handler.NewAuth(npy.Services.User, i)
	auth.RegisterRoutes(e.Group(""))
}

func onboardRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	onboard := handler.NewOnboard(npy.Services.Onboard)
	onboard.RegisterRoutes(e.Group("/onboard"), i)
}

func dashboardRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	dashboard := handler.NewDashboard(npy.Services, i)
	dashboard.RegisterRoutes(e.Group("/dashboard", middleware.JWTMiddleware()))
}

func userRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	user := handler.NewUser(npy.Services.User, i)
	user.RegisterRoutes(e.Group("/users", middleware.JWTMiddleware()))
}

func applicationRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	application := handler.NewApplication(npy.Services.Application, i)
	application.RegisterRoutes(e.Group("/applications", middleware.JWTMiddleware()))
}

func roleRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	roleHandler := handler.NewRole(i, npy.Services.Role)
	roleHandler.RegisterRoutes(e.Group("/roles", middleware.JWTMiddleware()))
}

func metadataRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	metadata := handler.NewMetadata(npy.Services.Metadata, i)
	metadata.RegisterRoutes(e.Group("/metadata", middleware.JWTMiddleware()))
}

func techStackRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	techStack := handler.NewTechStack(i, npy.Services.TechStack)
	techStack.RegisterRoutes(e.Group("/tech-stacks", middleware.JWTMiddleware()))
}

func gatewayRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	gateway := handler.NewGateway(npy.Services.Gateway, npy.Services.HealthChecker, i)
	gateway.RegisterRoutes(e.Group("/gateways", middleware.JWTMiddleware()))
}

func RegisterRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	loginRoutes(e, i, npy)
	onboardRoutes(e, i, npy)
	dashboardRoutes(e, i, npy)
	userRoutes(e, i, npy)
	applicationRoutes(e, i, npy)
	roleRoutes(e, i, npy)
	metadataRoutes(e, i, npy)
	techStackRoutes(e, i, npy)
	gatewayRoutes(e, i, npy)

	gateways, _ := npy.Services.Gateway.GetAll(context.Background())
	for _, gateway := range gateways {
		route := neployway.Route{
			AppID:     gateway.ApplicationID,
			Port:      gateway.Port,
			Domain:    gateway.Domain,
			Path:      gateway.Path,
			Subdomain: gateway.Subdomain,
		}
		if err := npy.Router.AddRoute(route); err != nil {
			logger.Error("Failed to add route: %v", err)
		}
	}

	// Use the router as a fallback handler for unmatched routes
	e.Any("/*", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		npy.Router.ServeHTTP(w, r)
	})))
}
