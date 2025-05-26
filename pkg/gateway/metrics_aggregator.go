package gateway

import (
	"context"
	"log"
	"os"
	"strings"
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
	m.mu.RLock()
	defer m.mu.RUnlock()

	for appID, collector := range m.collectors {
		metrics, err := collector.GetMetrics(30) // últimos 30 días
		if err != nil {
			log.Printf("ERROR: Failed to get metrics for app %s: %v", appID, err)
			continue
		}

		processedHours := make(map[string]bool)

		for _, metric := range metrics {
			timestamp, err := time.Parse("2006-01-02 15:00", metric.Hour)
			if err != nil {
				log.Printf("WARN: Invalid time format in metrics: %s", metric.Hour)
				continue
			}

			stat := model.ApplicationStat{
				ApplicationID: metric.ApplicationID,
				Requests:      metric.Requests,
				Errors:        metric.Errors,
				Date:          model.Date{Time: timestamp},
			}

			if err := m.appStatRepo.Insert(context.Background(), stat); err != nil {
				log.Printf("ERROR: Failed to save metrics for app %s: %v", appID, err)
				continue
			}

			// Marca la hora como procesada
			processedHours[metric.Hour] = true
		}

		// 🔥 Limpiar del log solo las horas procesadas
		removeProcessedLines(collector.metricsFile, processedHours)
	}
}

func removeProcessedLines(filePath string, processedHours map[string]bool) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("ERROR: reading metrics file to clean: %v", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	filtered := make([]string, 0)

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " - ")
		if len(parts) < 1 {
			continue
		}
		hourPart := parts[0]
		if !processedHours[hourPart] {
			filtered = append(filtered, line)
		}
	}

	err = os.WriteFile(filePath, []byte(strings.Join(filtered, "\n")), 0644)
	if err != nil {
		log.Printf("ERROR: writing cleaned metrics file: %v", err)
	}
}
