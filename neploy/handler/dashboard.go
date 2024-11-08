package handler

import (
	"context"
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
	service  service.Metadata
	sessions *session.Store
}

func NewDashboard(metadata service.Metadata, session *session.Store) *Dashboard {
	return &Dashboard{metadata, session}
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
			return
		}

		// parse the jwt token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Env.JWTSecret), nil
		}, nil)
		if err != nil {
			log.Err(err).Msg("error parsing token")
			return
		}

		// check if the token is valid
		if !token.Valid {
			log.Err(err).Msg("token is invalid")
			return
		}

		// get the role from the token
		role := token.Claims.(model.JWTClaims).Email

		admin := true
		if role != "henrrybrgt@gmail.com" {
			admin = false
		}

		teamName, err := d.service.GetTeamName(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting teamname")
			return
		}

		primaryColor, err := d.service.GetPrimaryColor(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting primary color")
			return
		}

		secondaryColor, err := d.service.GetSecondaryColor(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting secondary color")
		}

		logoUrl, err := d.service.GetTeamLogo(context.Background())
		if err != nil {
			log.Err(err).Msg("error getting logo")
			return
		}

		i.Render(w, r, "Dashboard/Index", gonertia.Props{
			"teamName":       teamName,
			"primaryColor":   primaryColor,
			"secondaryColor": secondaryColor,
			"logoUrl":        logoUrl,
			"admin":          admin,
		})
	}
}
