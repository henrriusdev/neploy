package middleware

import (
	"context"
	"neploy.dev/pkg/model"
	"net/http"
	"time"

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
				return c.Redirect(http.StatusSeeOther, "")
			}

			// Check if token exists
			if cookie.Value == "" {
				return c.Redirect(http.StatusSeeOther, "")
			}

			// Validate JWT token
			claims, valid, err := service.ValidateJWT(cookie.Value)
			if err != nil || !valid {
				return c.Redirect(http.StatusSeeOther, "")
			}

			// Store claims in context
			c.Set("claims", claims)
			return next(c)
		}
	}
}

func TraceMiddleware(traceService service.Trace) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("claims").(model.JWTClaims)
			if !ok {
				return next(c)
			}

			trace := &model.Trace{
				UserID:          claims.ID,
				Type:            "panel", // o dinámico según ruta
				Action:          c.Request().Method + " " + c.Path(),
				ActionTimestamp: model.NewDateNow(),
			}

			// Inyectar en contexto
			ctx := common.InjectTrace(c.Request().Context(), trace)
			c.SetRequest(c.Request().WithContext(ctx))

			err := next(c)

			// Guardar al final
			go traceService.Create(context.Background(), *trace)
			return err
		}
	}
}

func VisitorTraceMiddleware(visitorTraceService service.Visitor) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := int(time.Since(start).Milliseconds())

			visitorID := c.Request().Header.Get("X-Visitor-ID")
			if visitorID == "" {
				visitorID = c.RealIP()
			}

			appID := c.Request().Header.Get("X-Application-ID") // asegúrate de inyectarlo
			trace := model.VisitorTrace{
				ApplicationID: appID,
				VisitorID:     visitorID,
				PageVisited:   c.Path(),
				VisitDuration: duration,
				VisitedAt:     model.NewDateNow(),
			}

			go visitorTraceService.CreateTrace(c.Request().Context(), trace)
			return err
		}
	}
}
