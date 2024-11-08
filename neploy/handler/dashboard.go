package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/romsar/gonertia"
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
		i.Render(w, r, "Dashboard/Index", gonertia.Props{})
	}
}
