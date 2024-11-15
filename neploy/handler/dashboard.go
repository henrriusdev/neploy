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
	metadata service.Metadata
	app      service.Application
	sessions *session.Store
}

func NewDashboard(metadata service.Metadata, app service.Application, sessions *session.Store) *Dashboard {
	return &Dashboard{metadata, app, sessions}
}

func (d *Dashboard) RegisterRoutes(r fiber.Router, i *gonertia.Inertia) {
	r.Get("", adaptor.HTTPHandler(d.Index(i)))
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
			return
		}

		role := claims.Email

		admin := role == "henrrybrgt@gmail.com"

		teamName, err := d.metadata.GetTeamName(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting teamname")
			return
		}

		primaryColor, err := d.metadata.GetPrimaryColor(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting primary color")
			return
		}

		secondaryColor, err := d.metadata.GetSecondaryColor(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting secondary color")
		}

		logoUrl, err := d.metadata.GetTeamLogo(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting logo")
			return
		}

		healthyApps, _, err := d.app.GetHealthy(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting healthy apps")
			return
		}

		i.Render(w, r, "Dashboard/Index", gonertia.Props{
			"teamName":       teamName,
			"primaryColor":   primaryColor,
			"secondaryColor": secondaryColor,
			"logoUrl":        logoUrl,
			"admin":          admin,
			"health":         fmt.Sprintf("%d/%d", healthyApps, 4),
		})
	}
}
