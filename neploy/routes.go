package neploy

import (
	"github.com/labstack/echo/v4"
	"neploy.dev/neploy/handler"

	inertia "github.com/romsar/gonertia"
)

func loginRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	auth := handler.NewAuth(npy.Validator, npy.Services.User)
	auth.RegisterRoutes(e.Group(""), i)
}

func onboardRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	onboard := handler.NewOnboard(npy.Validator, npy.Services.Onboard)
	onboard.RegisterRoutes(e.Group("/onboard"), i)
}

func dashboardRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	dashboard := handler.NewDashboard(npy.Services)
	dashboard.RegisterRoutes(e.Group("/dashboard"), i)
}

func userRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	user := handler.NewUser(npy.Services.User)
	user.RegisterRoutes(e.Group("/users"), i)
}

func applicationRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	application := handler.NewApplication(npy.Services.Application)
	application.RegisterRoutes(e.Group("/applications"), i)
}

func RegisterRoutes(e *echo.Echo, i *inertia.Inertia, npy Neploy) {
	loginRoutes(e, i, npy)
	onboardRoutes(e, i, npy)
	dashboardRoutes(e, i, npy)
	userRoutes(e, i, npy)
	applicationRoutes(e, i, npy)
}
