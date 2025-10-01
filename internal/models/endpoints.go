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
	ID     int    `json:"id,omitempty"`
	UUID   string `json:"uuid,omitempty"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	// Basic Health Check
	EndpointPath string `json:"endpoint,omitempty"` // e.g., "/health"
	Timeout      int    `json:"timeout,omitempty"`  // Default: 30s
	Interval     int    `json:"interval,omitempty"` // Default: 5min
	// SSL Configuration
	CheckSSL bool `json:"check_ssl"`
	// Current Status
	Status               EndpointStatus `json:"status"`
	ResponseTime         int            `json:"response_time,omitempty"`
	ResponseMessage      string         `json:"response_message,omitempty"`
	ExpectedResponseCode int            `json:"expected_response_code,omitempty"`
	TimeoutSeconds       int            `json:"timeout_seconds,omitempty"`
	// Alert Groups (referência aos grupos de alerta)
	AlertGroupID *int `json:"alert_group_id,omitempty"`
	// Control Fields
	Enabled   bool      `json:"enabled"`
	LastCheck time.Time `json:"last_check,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
