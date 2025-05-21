package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type HealthChecker interface {
	Start(ctx context.Context)
	Stop()
	CheckGatewayHealth(ctx context.Context, gateway model.Gateway) error
}

type healthChecker struct {
	gatewayRepo *repository.Gateway
	appRepo     *repository.Application
	interval    time.Duration
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewHealthChecker(gatewayRepo *repository.Gateway, appRepo *repository.Application, interval time.Duration) HealthChecker {
	return &healthChecker{
		gatewayRepo: gatewayRepo,
		appRepo:     appRepo,
		interval:    interval,
		stopChan:    make(chan struct{}),
	}
}

func (h *healthChecker) Start(ctx context.Context) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		ticker := time.NewTicker(h.interval)
		defer ticker.Stop()

		for {
			select {
			case <-h.stopChan:
				return
			case <-ticker.C:
				h.checkAllGateways(ctx)
			}
		}
	}()
}

func (h *healthChecker) Stop() {
	close(h.stopChan)
	h.wg.Wait()
}

func (h *healthChecker) checkAllGateways(ctx context.Context) {
	gateways, err := h.gatewayRepo.GetAll(ctx)
	if err != nil {
		fmt.Printf("Error fetching gateways: %v\n", err)
		return
	}

	for _, gateway := range gateways {
		if err := h.CheckGatewayHealth(ctx, gateway); err != nil {
			fmt.Printf("Health check failed for gateway %s: %v\n", gateway.ID, err)
			// Update gateway status to unhealthy
			gateway.Status = "unhealthy"
			if err := h.gatewayRepo.Update(ctx, gateway); err != nil {
				fmt.Printf("Failed to update gateway status: %v\n", err)
			}
		} else {
			// Update gateway status to healthy if it was unhealthy before
			if gateway.Status != "healthy" {
				gateway.Status = "healthy"
				if err := h.gatewayRepo.Update(ctx, gateway); err != nil {
					fmt.Printf("Failed to update gateway status: %v\n", err)
				}
			}
		}
	}
}

func (h *healthChecker) CheckGatewayHealth(ctx context.Context, gateway model.Gateway) error {
	// Construct health check URL based on gateway configuration
	var healthCheckURL string = fmt.Sprintf("http://%s%s/health", gateway.Domain, gateway.Path)
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make health check request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthCheckURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}
