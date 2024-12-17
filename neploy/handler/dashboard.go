package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
	"github.com/romsar/gonertia"
	"github.com/rs/zerolog/log"
	"neploy.dev/config"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/service"
)

type Dashboard struct {
	services service.Services
	sessions *session.Store
}

func NewDashboard(services service.Services, sessions *session.Store) *Dashboard {
	return &Dashboard{
		services: services,
		sessions: sessions,
	}
}

func (d *Dashboard) RegisterRoutes(r fiber.Router, i *gonertia.Inertia) {
	r.Get("", adaptor.HTTPHandler(d.Index(i)))
	r.Get("/team", adaptor.HTTPHandler(d.Team(i)))
}

func (d *Dashboard) Index(i *gonertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the cookies
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Err(err).Msg("error getting token")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// parse the jwt token
		claims := &model.JWTClaims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Env.JWTSecret), nil
		})
		if err != nil {
			log.Err(err).Msg("error parsing token")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		role := claims.Email

		admin := role == "henrrybrgt@gmail.com"

		metadata, err := d.services.Metadata.Get(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting metadata")
			return
		}

		healthyApps, _, err := d.services.Application.GetHealthy(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting healthy apps")
			return
		}

		provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
		if err != nil {
			log.Err(err).Msg("error getting provider")
			return
		}

		user := model.UserResponse{
			Email:    claims.Email,
			Username: claims.Username,
			Name:     claims.Name,
			Provider: provider,
		}

		i.Render(w, r, "Dashboard/Index", gonertia.Props{
			"teamName": metadata.TeamName,
			"logoUrl":  metadata.LogoURL,
			"admin":    admin,
			"health":   fmt.Sprintf("%d/%d", healthyApps, 4),
			"user":     user,
		})
	}
}

func (d *Dashboard) Team(i *gonertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the cookies
		token, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		// get user data
		claims := &model.JWTClaims{}
		_, err = jwt.ParseWithClaims(token.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Env.JWTSecret), nil
		})
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		provider, err := d.services.User.GetProvider(context.Background(), claims.ID)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		user := model.UserResponse{
			Email:    claims.Email,
			Username: claims.Username,
			Name:     claims.Name,
			Provider: provider,
		}

		roles, err := d.services.Role.Get(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get metadata
		metadata, err := d.services.Metadata.Get(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		team, err := d.services.User.List(context.Background(), 15, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		i.Render(w, r, "Dashboard/Team", map[string]interface{}{
			"user":     user,
			"teamName": metadata.TeamName,
			"logoUrl":  metadata.LogoURL,
			"team":     team,
			"roles":    roles,
		})
	}
}
