package models

import (
	"time"
)

type EndpointStatus string

const (
	StatusOnline  EndpointStatus = "online"
	StatusOffline EndpointStatus = "offline"
	StatusUnknown EndpointStatus = "unknown"
)

type Endpoint struct {
	ID                   int            `json:"id,omitempty"`
	UUID                 string         `json:"uuid,omitempty"`
	Name                 string         `json:"name"`
	Domain               string         `json:"domain"`
	EndpointPath         string         `json:"endpoint,omitempty"` // e.g., "/health"
	Timeout              int            `json:"timeout,omitempty"`  // Default: 30s
	CheckSSL             bool           `json:"check_ssl,omitempty"`
	Status               EndpointStatus `json:"status"`
	ResponseTime         int            `json:"response_time,omitempty"`
	ResponseMessage      string         `json:"response_message,omitempty"`
	ExpectedResponseCode int            `json:"expected_response_code,omitempty"`
	ResponseStatusCode   int            `json:"response_code,omitempty"`
	TimeoutSeconds       int            `json:"timeout_seconds,omitempty"`
	AlertGroupID         *int           `json:"alert_group_id,omitempty"`
	Enabled              bool           `json:"enabled,omitempty"`
	LastCheck            time.Time      `json:"last_check,omitempty"`
	CreatedAt            time.Time      `json:"created_at,omitempty"`
	UpdatedAt            time.Time      `json:"updated_at,omitempty"`
}

// EndpointHealthHistory - Para manter histórico de checks
type EndpointHealthHistory struct {
	ID           int            `json:"id,omitempty"`
	EndPointID   int            `json:"endpoint_id"`
	Status       EndpointStatus `json:"status"`
	ResponseTime time.Duration  `json:"response_time,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	CheckedAt    time.Time      `json:"checked_at"`
}

type EndpointResponse struct {
	ExpectedResponseCode int    `json:"expected_response_code,omitempty"`
	ResponseStatusCode   int    `json:"response_code,omitempty"`
	ResponseTime         int    `json:"response_time,omitempty"`
	ResponseMessage      string `json:"response_message,omitempty"`
	TimeoutSeconds       int    `json:"timeout_seconds,omitempty"`
}
