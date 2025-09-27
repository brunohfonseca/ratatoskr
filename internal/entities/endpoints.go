package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EndpointStatus string

const (
	StatusOnline  EndpointStatus = "online"
	StatusOffline EndpointStatus = "offline"
	StatusUnknown EndpointStatus = "unknown"
)

type Endpoint struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name   string             `bson:"name" json:"name" index:"unique"`
	Domain string             `bson:"domain" json:"domain" index:""`
	Port   int                `bson:"port,omitempty" json:"port,omitempty"`

	// Basic Health Check
	Endpoint string        `bson:"endpoint,omitempty" json:"endpoint,omitempty"` // e.g., "/health"
	Timeout  time.Duration `bson:"timeout,omitempty" json:"timeout,omitempty"`   // Default: 30s
	Interval time.Duration `bson:"interval,omitempty" json:"interval,omitempty"` // Default: 5min

	// SSL Configuration
	CheckSSL bool `bson:"check_ssl,omitempty" json:"check_ssl,omitempty"`
	SSLData  struct {
		ExpirationDate time.Time `bson:"expiration_date,omitempty" json:"expiration_date,omitempty"`
		Expired        bool      `bson:"expired" json:"expired"`
		DaysLeft       int       `bson:"days_left" json:"days_left"`
		Issuer         string    `bson:"issuer,omitempty" json:"issuer,omitempty"`
	} `bson:"ssl_data,omitempty" json:"ssl_data,omitempty"`

	// Current Status
	Status       EndpointStatus `bson:"status" json:"status" index:""`
	ResponseTime time.Duration  `bson:"response_time,omitempty" json:"response_time,omitempty"`
	ErrorMessage string         `bson:"error_message,omitempty" json:"error_message,omitempty"`

	// Alert Groups (referência aos grupos de alerta)
	AlertGroupIDs []primitive.ObjectID `bson:"alert_group_ids,omitempty" json:"alert_group_ids,omitempty"`

	// Authentication
	Authentication interface{} `bson:"authentication,omitempty" json:"authentication,omitempty"`

	// Control Fields
	Enabled   bool      `bson:"enabled" json:"enabled" index:""`
	LastCheck time.Time `bson:"last_check,omitempty" json:"last_check,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// EndpointHealthHistory - Para manter histórico de checks
type EndpointHealthHistory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EndPointID   primitive.ObjectID `bson:"endpoint_id" json:"endpoint_id" index:""`
	Status       EndpointStatus     `bson:"status" json:"status" index:""`
	ResponseTime time.Duration      `bson:"response_time,omitempty" json:"response_time,omitempty"`
	ErrorMessage string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
	CheckedAt    time.Time          `bson:"checked_at" json:"checked_at" index:"desc" ttl:"120d"`
}
