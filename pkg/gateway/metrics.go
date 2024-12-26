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
}

func NewMetricsCollector(dataDir string) (*MetricsCollector, error) {
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
	}, nil
}

func (m *MetricsCollector) RecordRequest(timestamp time.Time, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Format the hour key: "2024-12-26 11:00"
	hourKey := timestamp.Format("2006-01-02 15:00")

	// Get current metrics for this hour
	metrics := m.hourlyMetrics[hourKey]

	// Update metrics
	metrics.requests++
	if isError {
		metrics.errors++
	}

	m.hourlyMetrics[hourKey] = metrics

	// Write to file
	m.writeMetrics(hourKey, metrics.requests, metrics.errors)
}

func (m *MetricsCollector) writeMetrics(hourKey string, requests, errors int) {
	// Format the line: "2024-12-26 11:00 - 100, 5"
	line := fmt.Sprintf("%s - %d, %d\n", hourKey, requests, errors)

	// Open file in append mode
	file, err := os.OpenFile(m.metricsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("ERROR: Failed to open metrics file: %v", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(line); err != nil {
		log.Printf("ERROR: Failed to write metrics: %v", err)
	}
}

func (m *MetricsCollector) GetMetrics(days int) ([]struct {
	Hour     string
	Requests int
	Errors   int
}, error,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Read the entire file
	content, err := os.ReadFile(m.metricsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics file: %v", err)
	}

	// Parse lines and aggregate data
	lines := strings.Split(string(content), "\n")
	metrics := make([]struct {
		Hour     string
		Requests int
		Errors   int
	}, 0)

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
		_, err := fmt.Sscanf(parts[1], "%d, %d", &requests, &errors)
		if err != nil {
			continue
		}

		// Parse the timestamp
		timestamp, err := time.Parse("2006-01-02 15:04", hourPart)
		if err != nil {
			continue
		}

		// Only include data after cutoff
		if timestamp.After(cutoff) {
			metrics = append(metrics, struct {
				Hour     string
				Requests int
				Errors   int
			}{
				Hour:     hourPart,
				Requests: requests,
				Errors:   errors,
			})
		}
	}

	return metrics, nil
}
