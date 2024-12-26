package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	inertia "github.com/romsar/gonertia"
	"github.com/rs/zerolog/log"
	"neploy.dev/config"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Dashboard struct {
	services service.Services
}

func NewDashboard(services service.Services) *Dashboard {
	return &Dashboard{
		services: services,
	}
}

func (d *Dashboard) RegisterRoutes(r *echo.Group, i *inertia.Inertia) {
	r.GET("", d.Index(i))
	r.GET("/team", d.Team(i))
	r.GET("/applications", d.Applications(i))
}

func (d *Dashboard) Index(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			log.Err(err).Msg("error getting token")
			return c.Redirect(http.StatusSeeOther, "/")
		}

		claims := &model.JWTClaims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Env.JWTSecret), nil
		})
		if err != nil {
			log.Err(err).Msg("error parsing token")
			return c.Redirect(http.StatusSeeOther, "/")
		}

		roles, err := d.services.Role.GetUserRoles(context.Background(), claims.ID)
		if err != nil {
			log.Err(err).Msg("error checking admin status")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		metadata, err := d.services.Metadata.Get(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting metadata")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		healthyApps, _, err := d.services.Application.GetHealthy(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting healthy apps")
			return err
		}

		provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
		if err != nil {
			log.Err(err).Msg("error getting provider")
			return err
		}

		user := model.UserResponse{
			Email:    claims.Email,
			Username: claims.Username,
			Name:     claims.Name,
			Provider: provider,
		}

		return i.Render(c.Response(), c.Request(), "Dashboard/Index", inertia.Props{
			"teamName": metadata.TeamName,
			"logoUrl":  metadata.LogoURL,
			"roles":    roles,
			"health":   fmt.Sprintf("%d/%d", healthyApps, 4),
			"user":     user,
		})
	}
}

func (d *Dashboard) Team(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.Redirect(http.StatusFound, "/auth/login")
		}

		claims := &model.JWTClaims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Env.JWTSecret), nil
		})
		if err != nil {
			return c.Redirect(http.StatusFound, "/auth/login")
		}

		provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
		if err != nil {
			return c.Redirect(http.StatusFound, "/auth/login")
		}

		user := model.UserResponse{
			Email:    claims.Email,
			Username: claims.Username,
			Name:     claims.Name,
			Provider: provider,
		}

		roles, err := d.services.Role.Get(context.Background())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		metadata, err := d.services.Metadata.Get(context.Background())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		listResponse, err := d.services.User.List(context.Background(), 15, 0)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return i.Render(c.Response(), c.Request(), "Dashboard/Team", inertia.Props{
			"user":     user,
			"teamName": metadata.TeamName,
			"logoUrl":  metadata.LogoURL,
			"team":     listResponse,
			"roles":    roles,
		})
	}
}

func (d *Dashboard) Applications(i *inertia.Inertia) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			log.Err(err).Msg("error getting token")
			return c.Redirect(http.StatusSeeOther, "/")
		}

		claims := &model.JWTClaims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Env.JWTSecret), nil
		})
		if err != nil {
			log.Err(err).Msg("error parsing token")
			return c.Redirect(http.StatusSeeOther, "/")
		}

		applications, err := d.services.Application.GetAll(c.Request().Context())
		if err != nil {
			log.Err(err).Msg("error getting applications")
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}

		metadata, err := d.services.Metadata.Get(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
		if err != nil {
			log.Err(err).Msg("error getting provider")
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}

		user := model.UserResponse{
			Email:    claims.Email,
			Username: claims.Username,
			Name:     claims.Name,
			Provider: provider,
		}

		props := inertia.Props{
			"user":         user,
			"teamName":     metadata.TeamName,
			"logoUrl":      metadata.LogoURL,
			"applications": applications,
		}

		if err := i.Render(c.Response(), c.Request(), "Dashboard/Applications", props); err != nil {
			log.Err(err).Msg("error rendering applications page")
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}

		return nil
	}
}
