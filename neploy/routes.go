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

func RegisterRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	loginRoutes(e, i, npy)
	onboardRoutes(e, i, npy)
	dashboardRoutes(e, i, npy)
	userRoutes(e, i, npy)
	applicationRoutes(e, i, npy)
	roleRoutes(e, i, npy)

	gateways, _ := npy.Services.Gateway.GetAll(context.Background())
	router := neployway.NewRouter(npy.Repositories.ApplicationStat)
	for _, gateway := range gateways {
		route := neployway.Route{
			AppID:     gateway.ApplicationID,
			Port:      gateway.Port,
			Domain:    gateway.Domain,
			Path:      gateway.Path,
			Subdomain: gateway.Subdomain,
		}
		if err := router.AddRoute(route); err != nil {
			logger.Error("Failed to add route: %v", err)
		}
	}

	// Use the router as a fallback handler for unmatched routes
	e.Any("/*", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	})))
}
