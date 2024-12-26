package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/romsar/gonertia"
	"github.com/rs/zerolog/log"
	"neploy.dev/neploy/validation"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Onboard struct {
	validator validation.XValidator
	service   service.Onboard
}

func NewOnboard(validator validation.XValidator, service service.Onboard) *Onboard {
	return &Onboard{
		validator: validator,
		service:   service,
	}
}

func (o *Onboard) RegisterRoutes(r *echo.Group, i *gonertia.Inertia) {
	r.POST("", o.Initiate)
}

func (o *Onboard) Initiate(c echo.Context) error {
	var req model.OnboardRequest
	if err := c.Bind(&req); err != nil {
		log.Err(err).Msg("error parsing request")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request",
		})
	}

	oauthID, err := c.Cookie("oauth_id")
	if err != nil || oauthID.Value == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	// delete oauth_id cookie
	cookie := new(http.Cookie)
	cookie.Name = "oauth_id"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	req.OauthID = oauthID.Value

	if err := o.service.Initiate(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Onboarding initiated",
	})
}
