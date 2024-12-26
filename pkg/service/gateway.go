package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	neployway "neploy.dev/pkg/gateway"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Gateway interface {
	Create(ctx context.Context, gateway model.Gateway) error
	Update(ctx context.Context, gateway model.Gateway) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (model.Gateway, error)
	ListByApp(ctx context.Context, appID string) ([]model.Gateway, error)
	AddRoute(ctx context.Context, appID, domain, subdomain, path string) error
	RemoveRoute(ctx context.Context, appID, domain string) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type gateway struct {
	router *neployway.Router
	repo   repository.Gateway
	mu     sync.RWMutex
}

func NewGatewayService(repo repository.Gateway) Gateway {
	return &gateway{
		router: neployway.NewRouter(),
		repo:   repo,
	}
}

func (s *gateway) Create(ctx context.Context, gateway model.Gateway) error {
	gateway.ID = uuid.New().String()
	if err := s.repo.Insert(ctx, gateway); err != nil {
		return err
	}

	// Add route to router
	return s.AddRoute(ctx, gateway.ApplicationID, gateway.Domain, gateway.Subdomain, gateway.Path)
}

func (s *gateway) Update(ctx context.Context, gateway model.Gateway) error {
	oldGateway, err := s.Get(ctx, gateway.ID)
	if err != nil {
		return err
	}

	// Remove old route
	if err := s.RemoveRoute(ctx, oldGateway.ApplicationID, oldGateway.Domain); err != nil {
		return err
	}

	// Add new route
	if err := s.AddRoute(ctx, gateway.ApplicationID, gateway.Domain, gateway.Subdomain, gateway.Path); err != nil {
		return err
	}

	return s.repo.Update(ctx, gateway)
}

func (s *gateway) Delete(ctx context.Context, id string) error {
	gateway, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Remove route
	if err := s.RemoveRoute(ctx, gateway.ApplicationID, gateway.Domain); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *gateway) Get(ctx context.Context, id string) (model.Gateway, error) {
	gateway, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Gateway{}, errors.Wrap(err, "failed to get gateway")
	}
	return gateway, nil
}

func (s *gateway) ListByApp(ctx context.Context, appID string) ([]model.Gateway, error) {
	gateways, err := s.repo.GetByApplicationID(ctx, appID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list gateways")
	}
	return gateways, nil
}

func (s *gateway) AddRoute(ctx context.Context, appID, domain, subdomain, path string) error {
	// Get application details
	app, err := s.repo.GetApplication(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to get application: %v", err)
	}

	// Create route
	route := neployway.Route{
		AppID:     appID,
		Port:      app.Port,
		Domain:    domain,
		Subdomain: subdomain,
		Path:      path,
	}

	// Add route to router
	if err := s.router.AddRoute(route); err != nil {
		return fmt.Errorf("failed to add route: %v", err)
	}

	return nil
}

func (s *gateway) RemoveRoute(ctx context.Context, appID, domain string) error {
	// Get route from database
	route, err := s.repo.GetGateway(ctx, appID, domain)
	if err != nil {
		return fmt.Errorf("failed to get route: %v", err)
	}

	// Remove from router
	s.router.RemoveRoute(neployway.Route{
		AppID:     route.ApplicationID,
		Domain:    route.Domain,
		Subdomain: route.Subdomain,
		Path:      route.Path,
	})

	return nil
}

func (s *gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
