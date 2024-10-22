package common

import (
	"time"

	"neploy.dev/pkg/model"
)

func FormatDateRange(startDate, endDate time.Time) model.DateRange {
	return model.DateRange{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}
}
