package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"neploy.dev/pkg/common"
	"neploy.dev/pkg/service"
)

type AuthConfig struct {
	SessionStore *session.Store
}

func OnboardingMiddleware(service service.Onboard) fiber.Handler {
	onboardPath := "/onboard" // Change this to match your onboarding path

	return func(c *fiber.Ctx) error {
		// Skip middleware for non-GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		if common.AcceptedRoutesForOnboarding(c.Path()) {
			return c.Next()
		}

		// Check if onboarding is completed
		isDone, err := service.Done(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check onboarding status",
			})
		}

		// If onboarding is not done, handle the redirect
		if !isDone {
			if c.Path() == onboardPath {
				return c.Next()
			}
			// For regular requests, do a normal redirect
			return c.Redirect(onboardPath)
		}

		// If onboarding is done and user tries to access onboard page,
		// handle the redirect to home (optional)

		if isDone && c.Path() == onboardPath {
			return c.Redirect("/")
		}

		return c.Next()
	}
}

// JWTMiddleware is a middleware that checks if the user is authenticated
func SessionMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Next()
		}

		// Store session in locals for easy access
		c.Locals("session", sess)
		return c.Next()
	}
}

func AuthMiddleware(config AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := config.SessionStore.Get(c)
		if err != nil {
			return c.Redirect("/login")
		}

		auth := sess.Get("authenticated")
		if auth == nil || !auth.(bool) {
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
