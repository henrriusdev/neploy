package gateway

import (
	"context"
	"log"
	"sync"
	"time"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type MetricsAggregator struct {
	metricsCollector *MetricsCollector
	appStatRepo      repository.ApplicationStat
	stopChan         chan struct{}
	wg               sync.WaitGroup
}

func NewMetricsAggregator(metricsCollector *MetricsCollector, appStatRepo repository.ApplicationStat) *MetricsAggregator {
	return &MetricsAggregator{
		metricsCollector: metricsCollector,
		appStatRepo:      appStatRepo,
		stopChan:         make(chan struct{}),
	}
}

func (m *MetricsAggregator) Start() {
	m.wg.Add(1)
	go m.run()
}

func (m *MetricsAggregator) Stop() {
	close(m.stopChan)
	m.wg.Wait()
}

func (m *MetricsAggregator) run() {
	defer m.wg.Done()

	// Calculate time until next hour
	now := time.Now()
	nextHour := now.Truncate(time.Hour).Add(time.Hour)
	initialDelay := nextHour.Sub(now)

	// Wait until the next hour before starting
	timer := time.NewTimer(initialDelay)
	defer timer.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-timer.C:
			// Aggregate and save metrics for the previous hour
			m.aggregateAndSaveMetrics()

			// Reset timer for the next hour
			timer.Reset(time.Hour)
		}
	}
}

func (m *MetricsAggregator) aggregateAndSaveMetrics() {
	// Get metrics for the last hour
	endTime := time.Now()
	startTime := endTime.Add(-time.Hour)

	metrics, err := m.metricsCollector.GetMetrics(1) // Get last 24 hours of metrics
	if err != nil {
		log.Printf("ERROR: Failed to get metrics: %v", err)
		return
	}

	// Find metrics for the previous hour
	var lastHourMetrics LastHourMetrics

	targetHour := startTime.Format("2006-01-02 15:04")
	for _, m := range metrics {
		if m.Hour == targetHour {
			lastHourMetrics = m
			break
		}
	}

	// Create application stat
	stat := model.ApplicationStat{
		ApplicationID: "default", // TODO: Update when we add per-app metrics
		Date:          model.Date{Time: startTime},
		Requests:      lastHourMetrics.Requests,
		Errors:        lastHourMetrics.Errors,
		// TODO: Add these fields when implemented
		// AverageResponseTime: 0,
		// DataTransfered:      0,
		// UniqueVisitors:      0,
	}

	// Save to database
	ctx := context.Background()
	if err := m.appStatRepo.Insert(ctx, stat); err != nil {
		log.Printf("ERROR: Failed to save application stats: %v", err)
		return
	}

	log.Printf("INFO: Successfully saved application stats for hour %s", targetHour)
}
