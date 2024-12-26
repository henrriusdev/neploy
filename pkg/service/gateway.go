package service

import (
	"context"
	"fmt"
	"net/http"

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
	AddRoute(ctx context.Context, gateway model.Gateway) error
	RemoveRoute(ctx context.Context, gateway model.Gateway) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type gateway struct {
	router  *neployway.Router
	repo    repository.Gateway
	appRepo repository.Application
}

func NewGatewayService(repo repository.Gateway, appRepo repository.Application, statRepo repository.ApplicationStat) Gateway {
	return &gateway{
		router:  neployway.NewRouter(statRepo),
		repo:    repo,
		appRepo: appRepo,
	}
}

func (s *gateway) validateGateway(ctx context.Context, gateway model.Gateway) error {
	// Validate required fields
	if gateway.Name == "" {
		return errors.New("gateway name is required")
	}
	if gateway.ApplicationID == "" {
		return errors.New("application ID is required")
	}
	if gateway.EndpointType == "" {
		return errors.New("endpoint type is required")
	}
	if gateway.Domain == "" {
		return errors.New("domain is required")
	}

	// Validate endpoint type specific fields
	switch gateway.EndpointType {
	case "subdomain":
		if gateway.Subdomain == "" {
			return errors.New("subdomain is required for subdomain endpoint type")
		}
	case "path":
		if gateway.Path == "" {
			return errors.New("path is required for path endpoint type")
		}
	default:
		return errors.New("invalid endpoint type")
	}

	// Check if application exists
	_, err := s.appRepo.GetByID(ctx, gateway.ApplicationID)
	if err != nil {
		return errors.Wrap(err, "application not found")
	}

	return nil
}

func (s *gateway) Create(ctx context.Context, gateway model.Gateway) error {
	if err := s.validateGateway(ctx, gateway); err != nil {
		return err
	}

	gateway.ID = uuid.New().String()
	gateway.Status = "active"

	if err := s.repo.Insert(ctx, gateway); err != nil {
		return err
	}

	return s.AddRoute(ctx, gateway)
}

func (s *gateway) Update(ctx context.Context, gateway model.Gateway) error {
	if err := s.validateGateway(ctx, gateway); err != nil {
		return err
	}

	oldGateway, err := s.Get(ctx, gateway.ID)
	if err != nil {
		return err
	}

	// Remove old route
	if err := s.RemoveRoute(ctx, oldGateway); err != nil {
		return err
	}

	// Add new route
	if err := s.AddRoute(ctx, gateway); err != nil {
		return err
	}

	return s.repo.Update(ctx, gateway)
}

func (s *gateway) Delete(ctx context.Context, id string) error {
	gateway, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := s.RemoveRoute(ctx, gateway); err != nil {
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

func (s *gateway) AddRoute(ctx context.Context, gateway model.Gateway) error {
	route := neployway.Route{
		AppID:     gateway.ApplicationID,
		Port:      gateway.Port,
		Domain:    gateway.Domain,
		Subdomain: gateway.Subdomain,
		Path:      gateway.Path,
	}

	if err := s.router.AddRoute(route); err != nil {
		gateway.Status = "error"
		s.repo.Update(ctx, gateway)
		return fmt.Errorf("failed to add route: %v", err)
	}

	gateway.Status = "active"
	return s.repo.Update(ctx, gateway)
}

func (s *gateway) RemoveRoute(ctx context.Context, gateway model.Gateway) error {
	route := neployway.Route{
		AppID:     gateway.ApplicationID,
		Domain:    gateway.Domain,
		Subdomain: gateway.Subdomain,
		Path:      gateway.Path,
	}

	s.router.RemoveRoute(route)
	return nil
}

func (s *gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
