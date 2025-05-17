package gateway

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type MetricsCollector struct {
	mu            sync.Mutex
	metricsFile   string
	hourlyMetrics map[string]struct {
		requests int
		errors   int
	}
	applicationID string
}

func NewMetricsCollector(dataDir string, applicationID string) (*MetricsCollector, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create metrics directory: %v", err)
	}

	metricsFile := filepath.Join(dataDir, "gateway_metrics.log")

	return &MetricsCollector{
		metricsFile: metricsFile,
		hourlyMetrics: make(map[string]struct {
			requests int
			errors   int
		}),
		applicationID: applicationID,
	}, nil
}

func (m *MetricsCollector) RecordRequest(timestamp time.Time, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	hourKey := timestamp.Format("2006-01-02 15:00")

	metrics := m.hourlyMetrics[hourKey]
	metrics.requests++
	if isError {
		metrics.errors++
	}

	m.hourlyMetrics[hourKey] = metrics
	m.writeMetrics(hourKey, metrics.requests, metrics.errors)
}

func (m *MetricsCollector) writeMetrics(hourKey string, requests, errors int) {
	lineToWrite := fmt.Sprintf("%s - %d, %d = %s", hourKey, requests, errors, m.applicationID)

	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	// Leer el archivo actual
	content, err := os.ReadFile(m.metricsFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("ERROR: reading metrics file: %v", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	updated := false

	for i, line := range lines {
		if strings.HasPrefix(line, hourKey+" - ") {
			lines[i] = lineToWrite
			updated = true
			break
		}
	}

	if !updated {
		lines = append(lines, lineToWrite)
	}

	finalContent := strings.Join(lines, "\n")
	err = os.WriteFile(m.metricsFile, []byte(finalContent), 0644)
	if err != nil {
		log.Printf("ERROR: writing metrics file: %v", err)
	}
}

func (m *MetricsCollector) GetMetrics(days int) ([]LastHourMetrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Read the entire file
	content, err := os.ReadFile(m.metricsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics file: %v", err)
	}

	// Parse lines and aggregate data
	lines := strings.Split(string(content), "\n")
	metrics := make([]LastHourMetrics, 0)

	cutoff := time.Now().AddDate(0, 0, -days)

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse line: "2024-12-26 11:00 - 100, 5"
		parts := strings.Split(line, " - ")
		if len(parts) != 2 {
			continue
		}

		hourPart := parts[0]
		var requests, errors int
		var applicationID string
		_, err := fmt.Sscanf(parts[1], "%d, %d = %s", &requests, &errors, &applicationID)
		if err != nil {
			continue
		}

		// Parse the timestamp
		timestamp, err := time.Parse("2006-01-02 15:00", hourPart)
		if err != nil {
			continue
		}

		// Only include data after cutoff
		if timestamp.After(cutoff) {
			metrics = append(metrics, LastHourMetrics{
				Hour:          hourPart,
				Requests:      requests,
				Errors:        errors,
				ApplicationID: applicationID,
			})
		}
	}

	return metrics, nil
}
