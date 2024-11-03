package common

import (
	"strings"
	"time"

	"neploy.dev/pkg/model"
)

func FormatDateRange(startDate, endDate time.Time) model.DateRange {
	return model.DateRange{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}
}

func AcceptedRoutesForOnboarding(path string) bool {
	return strings.HasPrefix(path, "/build/assets/") || strings.HasPrefix(path, "/auth")
}
