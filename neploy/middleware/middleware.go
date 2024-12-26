package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"neploy.dev/pkg/common"
	"neploy.dev/pkg/service"
)

func OnboardingMiddleware(service service.Onboard) echo.MiddlewareFunc {
	onboardPath := "/onboard" // Change this to match your onboarding path

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip middleware for non-GET requests
			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			if common.AcceptedRoutesForOnboarding(c.Path()) {
				return next(c)
			}

			// Check if onboarding is completed
			isDone, err := service.Done(c.Request().Context())
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "Failed to check onboarding status",
				})
			}

			// If onboarding is not done, handle the redirect
			if !isDone {
				if c.Path() == onboardPath {
					return next(c)
				}
				// For regular requests, do a normal redirect
				return c.Redirect(http.StatusTemporaryRedirect, onboardPath)
			}

			// If onboarding is done and user tries to access onboard page,
			// handle the redirect to home
			if isDone && c.Path() == onboardPath {
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}

			return next(c)
		}
	}
}

// JWTMiddleware is a middleware that checks if the user is authenticated
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from cookie
			cookie, err := c.Cookie("token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authentication token")
			}

			// Check if token exists
			if cookie.Value == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authentication token")
			}

			// Validate JWT token
			claims, valid, err := service.ValidateJWT(cookie.Value)
			if err != nil || !valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication token")
			}

			// Store claims in context
			c.Set("claims", claims)
			return next(c)
		}
	}
}

// AuthMiddleware checks if the user is authenticated via JWT
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			if user == nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			return next(c)
		}
	}
}
