package gateway

type LastHourMetrics struct {
	Hour          string
	Requests      int
	Errors        int
	ApplicationID string
}
