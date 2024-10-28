package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"neploy.dev/pkg/service"
)

func OnboardingMiddleware(service service.Onboard) fiber.Handler {
	onboardPath := "/onboard" // Change this to match your onboarding path

	return func(c *fiber.Ctx) error {
		// Skip middleware for non-GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		if strings.HasPrefix(c.Path(), "/build/assets/") {
			return c.Next()
		}

		// Check if onboarding is completed
		isDone, step, err := service.Done(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check onboarding status",
			})
		}

		if strings.Contains(c.OriginalURL(), onboardPath) || strings.Contains(c.OriginalURL(), fmt.Sprintf("?step=%d", step)) {
			return c.Next()
		}

		// If onboarding is not done, handle the redirect
		if !isDone {
			// Check if it's an Inertia request
			if c.Get("X-Inertia") == "true" {
				// For Inertia requests, return a 409 Conflict with the redirect location
				c.Set("X-Inertia-Location", onboardPath)
				return c.SendStatus(fiber.StatusConflict)
			}

			// Only append the step query param if it's not already present
			if !strings.Contains(c.OriginalURL(), fmt.Sprintf("step=%d", step)) {
				onboardPath = fmt.Sprintf("%s?step=%d", onboardPath, step)
			}
			return c.Redirect(onboardPath)
		}

		// If onboarding is done and user tries to access onboard page,
		// handle the redirect to home (optional)
		if isDone && c.Path() == onboardPath {
			if c.Get("X-Inertia") == "true" {
				c.Set("X-Inertia-Location", "/")
				return c.SendStatus(fiber.StatusConflict)
			}
			return c.Redirect("/")
		}

		return c.Next()
	}
}

// JWTMiddleware is a middleware that checks if the user is authenticated
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip middleware for non-GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		if strings.HasPrefix(c.Path(), "/build/assets/") {
			return c.Next()
		}

		// Check if the user is authenticated
		if c.Locals("user") == nil {
			if c.Get("X-Inertia") == "true" {
				c.Set("X-Inertia-Location", "/login")
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
