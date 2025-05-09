package neploy

import (
	"context"
	"fmt"
	"neploy.dev/pkg/repository/filters"
	"net/http"
	"strconv"

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
	dashboard.RegisterRoutes(e.Group("/dashboard", middleware.JWTMiddleware(), middleware.TraceMiddleware(npy.Services.Trace)))
}

func userRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	user := handler.NewUser(npy.Services.User, npy.Services.Metadata, i)
	user.RegisterRoutes(e.Group("/users", middleware.JWTMiddleware(), middleware.TraceMiddleware(npy.Services.Trace)))
}

func applicationRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	application := handler.NewApplication(npy.Services.Application, i)
	application.RegisterRoutes(e.Group("/applications", middleware.JWTMiddleware(), middleware.TraceMiddleware(npy.Services.Trace)))
}

func roleRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	roleHandler := handler.NewRole(i, npy.Services.Role)
	roleHandler.RegisterRoutes(e.Group("/roles", middleware.JWTMiddleware(), middleware.TraceMiddleware(npy.Services.Trace)))
}

func metadataRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	metadata := handler.NewMetadata(npy.Services.Metadata, i)
	metadata.RegisterRoutes(e.Group("/metadata", middleware.JWTMiddleware(), middleware.TraceMiddleware(npy.Services.Trace)))
}

func techStackRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	techStack := handler.NewTechStack(i, npy.Services.TechStack)
	techStack.RegisterRoutes(e.Group("/tech-stacks", middleware.JWTMiddleware(), middleware.TraceMiddleware(npy.Services.Trace)))
}

func gatewayRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	gateway := handler.NewGateway(npy.Services.Gateway, npy.Services.HealthChecker, i)
	gateway.RegisterRoutes(e.Group("/gateways", middleware.JWTMiddleware(), middleware.TraceMiddleware(npy.Services.Trace)))
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
		// 🔄 Traer todas las versiones de esta app
		versions, err := npy.Repositories.ApplicationVersion.GetAll(context.Background(), filters.IsSelectFilter("application_id", gateway.ApplicationID))
		if err != nil {
			logger.Error("Failed to get versions for app %s: %v", gateway.ApplicationID, err)
			continue
		}

		for i, v := range versions {
			// Ruta con versión incluida: /vX.Y.Z/app
			versionedPath := fmt.Sprintf("/%s%s", v.VersionTag, gateway.EndpointURL)

			port := 4001
			port -= i

			route := neployway.Route{
				AppID:     gateway.ApplicationID,
				Port:      strconv.Itoa(int(port)),
				Domain:    gateway.Domain,
				Path:      versionedPath,
				Subdomain: gateway.Subdomain,
			}

			println("Registering default route:", route.Path, route.Port)

			if err := npy.Router.AddRoute(route); err != nil {
				logger.Error("Failed to add route: %v", err)
			}
		}

		// ➕ Registrar ruta sin versión para la versión por defecto
		route := neployway.Route{
			AppID:     gateway.ApplicationID,
			Port:      gateway.Port,
			Domain:    gateway.Domain,
			Path:      gateway.EndpointURL, // sin versión
			Subdomain: gateway.Subdomain,
		}

		println("Registering default route:", route.Path)

		if err := npy.Router.AddRoute(route); err != nil {
			logger.Error("Failed to add default route: %v", err)
		}
	}

	// Use the router as a fallback handler for unmatched routes
	e.Any("/*", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		npy.Router.ServeHTTP(w, r)
	})))
}
