package handler

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	inertia "github.com/romsar/gonertia"
	"log"
	"neploy.dev/config"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
	"net/http"
	"slices"
)

type Dashboard struct {
	services service.Services
	i        *inertia.Inertia
}

func NewDashboard(services service.Services, i *inertia.Inertia) *Dashboard {
	return &Dashboard{
		services: services,
		i:        i,
	}
}

func (d *Dashboard) RegisterRoutes(r *echo.Group) {
	r.GET("", d.Index)
	r.GET("/team", d.Team)
	r.GET("/applications", d.Applications)
	r.GET("/applications/:id", d.ApplicationView)
	r.GET("/gateways", d.Gateways)
	r.GET("/settings", d.Config)
}

func (d *Dashboard) Index(c echo.Context) error {
	claims, ok := c.Get("claims").(model.JWTClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	roles, err := d.services.Role.GetUserRoles(context.Background(), claims.ID)
	if err != nil {
		logger.Error("error checking admin status: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	metadata, err := d.services.Metadata.Get(context.Background())
	if err != nil {
		logger.Error("error getting metadata: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	healthyApps, _, err := d.services.Application.GetHealthy(context.Background())
	if err != nil {
		log.Println("error retrieving app health:", err)
		// manejar fallback si hace falta
	}

	provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
	if err != nil {
		logger.Error("error getting provider: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user := model.UserResponse{
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
		Roles:    claims.RolesLower,
		Provider: provider,
	}

	return d.i.Render(c.Response(), c.Request(), "Dashboard/Index", inertia.Props{
		"teamName": metadata.TeamName,
		"logoUrl":  metadata.LogoURL,
		"roles":    roles,
		"health":   fmt.Sprintf("%d/%d", healthyApps, 4),
		"user":     user,
	})
}

func (d *Dashboard) Team(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		return c.Redirect(http.StatusFound, "/auth/login")
	}

	claims := &model.JWTClaims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Env.JWTSecret), nil
	})
	if err != nil {
		logger.Error("error parsing token: %v", err)
		return c.Redirect(http.StatusFound, "/auth/login")
	}

	provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
	if err != nil {
		logger.Error("error getting provider: %v", err)
		return c.Redirect(http.StatusFound, "/auth/login")
	}

	user := model.UserResponse{
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
		Roles:    claims.RolesLower,
		Provider: provider,
	}

	roles, err := d.services.Role.Get(context.Background())
	if err != nil {
		logger.Error("error getting roles: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	metadata, err := d.services.Metadata.Get(context.Background())
	if err != nil {
		logger.Error("error getting metadata: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	listResponse, err := d.services.User.List(context.Background(), 15, 0)
	if err != nil {
		logger.Error("error listing users: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return d.i.Render(c.Response(), c.Request(), "Dashboard/Index", inertia.Props{
		"user":     user,
		"teamName": metadata.TeamName,
		"logoUrl":  metadata.LogoURL,
		"team":     listResponse,
		"roles":    roles,
	})
}

func (d *Dashboard) Applications(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		logger.Error("error getting token: %v", err)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	claims := &model.JWTClaims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Env.JWTSecret), nil
	})
	if err != nil {
		logger.Error("error parsing token: %v", err)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	applications, err := d.services.Application.GetAll(c.Request().Context(), claims.ID)
	if err != nil {
		logger.Error("error getting applications: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	metadata, err := d.services.Metadata.Get(c.Request().Context())
	if err != nil {
		logger.Error("error getting metadata: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
	if err != nil {
		logger.Error("error getting provider: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	user := model.UserResponse{
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
		Roles:    claims.RolesLower,
		Provider: provider,
	}

	props := inertia.Props{
		"user":         user,
		"teamName":     metadata.TeamName,
		"logoUrl":      metadata.LogoURL,
		"applications": applications,
	}

	if err := d.i.Render(c.Response(), c.Request(), "Dashboard/Index", props); err != nil {
		logger.Error("error rendering applications page: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return nil
}

func (d *Dashboard) ApplicationView(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		logger.Error("error getting token: %v", err)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	claims := &model.JWTClaims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Env.JWTSecret), nil
	})
	if err != nil {
		logger.Error("error parsing token: %v", err)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	application, err := d.services.Application.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		logger.Error("error getting applications: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	metadata, err := d.services.Metadata.Get(c.Request().Context())
	if err != nil {
		logger.Error("error getting metadata: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
	if err != nil {
		logger.Error("error getting provider: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	user := model.UserResponse{
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
		Roles:    claims.RolesLower,
		Provider: provider,
	}

	props := inertia.Props{
		"user":        user,
		"teamName":    metadata.TeamName,
		"logoUrl":     metadata.LogoURL,
		"application": application,
	}

	if err := d.i.Render(c.Response(), c.Request(), "Dashboard/Index", props); err != nil {
		logger.Error("error rendering applications page: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return nil
}

func (d *Dashboard) Gateways(c echo.Context) error {
	claims, ok := c.Get("claims").(model.JWTClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	gateways, err := d.services.Gateway.GetAll(c.Request().Context())
	if err != nil {
		logger.Error("error getting gateways: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	metadata, err := d.services.Metadata.Get(c.Request().Context())
	if err != nil {
		logger.Error("error getting metadata: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	provider, err := d.services.User.GetProvider(c.Request().Context(), claims.ID)
	if err != nil {
		logger.Error("error getting provider: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	conf, err := d.services.Gateway.GetConfig(c.Request().Context())
	if err != nil {
		logger.Error("error getting gateway config: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user := model.UserResponse{
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
		Roles:    claims.RolesLower,
		Provider: provider,
	}

	return d.i.Render(c.Response(), c.Request(), "Dashboard/Index", inertia.Props{
		"user":     user,
		"teamName": metadata.TeamName,
		"logoUrl":  metadata.LogoURL,
		"gateways": gateways,
		"config":   conf,
	})
}

func (d *Dashboard) Config(c echo.Context) error {
	claims, ok := c.Get("claims").(model.JWTClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if !slices.Contains(claims.RolesLower, "administrator") && !slices.Contains(claims.RolesLower, "settings") {
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}

	metadata, err := d.services.Metadata.Get(c.Request().Context())
	if err != nil {
		logger.Error("error getting metadata: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	provider, err := d.services.User.GetProvider(c.Request().Context(), claims.ID)
	if err != nil {
		logger.Error("error getting provider: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	user := model.UserResponse{
		Email:    claims.Email,
		Username: claims.Username,
		Name:     claims.Name,
		Roles:    claims.RolesLower,
		Provider: provider,
	}

	roles, err := d.services.Role.Get(c.Request().Context())
	if err != nil {
		logger.Error("error getting roles: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	techStacks, err := d.services.TechStack.GetAll(c.Request().Context())
	if err != nil {
		logger.Error("error getting tech stacks: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	traces, err := d.services.Trace.GetAll(c.Request().Context())
	if err != nil {
		logger.Error("error getting traces: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return d.i.Render(c.Response(), c.Request(), "Dashboard/Index", inertia.Props{
		"user":       user,
		"teamName":   metadata.TeamName,
		"logoUrl":    metadata.LogoURL,
		"language":   metadata.Language,
		"roles":      roles,
		"techStacks": techStacks,
		"traces":     traces,
	})
}
