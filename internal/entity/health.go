package entity

type Health string

const (
	HealthStatusOK     Health = "OK"
	HealthStatusFailed Health = "Failed"
)

type HealthCheck struct {
	Database Health `json:"database"`
}
