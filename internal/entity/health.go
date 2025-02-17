package entity

type Health string

const (
	HealthStatusOK Health = "OK"
	HealthStatusKO Health = "KO"
)

type HealthCheck struct {
	Database Health `json:"database"`
}
