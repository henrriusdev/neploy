package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"neploy.dev/config"
	"neploy.dev/pkg/common"
	"neploy.dev/pkg/model"
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
			if c.Path() == "/login" {
				return next(c)
			}
			// Get token from cookies
			token, err := c.Cookie("token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "missing authorization token",
				})
			}

			// Check if token is valid
			if token.Value == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "missing authorization token",
				})
			}

			tokenString := token.Value

			// Parse token
			jwt, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Replace this with your actual secret key
				return []byte(config.Env.JWTSecret), nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error":   "invalid token",
					"message": err.Error(),
				})
			}

			if claims, ok := jwt.Claims.(*model.JWTClaims); ok && jwt.Valid {
				// Store user info in context
				c.Set("claims", claims)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid token claims",
			})
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
