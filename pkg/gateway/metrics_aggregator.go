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
	collectors  map[string]*MetricsCollector // Map of appID to collector
	mu          sync.RWMutex
	appStatRepo *repository.ApplicationStat
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewMetricsAggregator(metricsCollector *MetricsCollector, appStatRepo *repository.ApplicationStat) *MetricsAggregator {
	collectors := make(map[string]*MetricsCollector)
	if metricsCollector != nil {
		collectors[metricsCollector.applicationID] = metricsCollector
	}

	return &MetricsAggregator{
		collectors:  collectors,
		appStatRepo: appStatRepo,
		stopChan:    make(chan struct{}),
	}
}

func (m *MetricsAggregator) AddCollector(collector *MetricsCollector) {
	if collector == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.collectors[collector.applicationID] = collector
}

func (m *MetricsAggregator) RemoveCollector(appID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.collectors, appID)
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

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Aggregate metrics for each collector
	for appID, collector := range m.collectors {
		metrics, err := collector.GetMetrics(1) // Get last 24 hours of metrics
		if err != nil {
			log.Printf("ERROR: Failed to get metrics for app %s: %v", appID, err)
			continue
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

		stat := model.ApplicationStat{
			ApplicationID: appID,
			Requests:      lastHourMetrics.Requests,
			Errors:        lastHourMetrics.Errors,
			Date:          model.Date{Time: startTime},
		}

		if err := m.appStatRepo.Insert(context.Background(), stat); err != nil {
			log.Printf("ERROR: Failed to save metrics for app %s: %v", appID, err)
		}
	}
}
