package entities

import "time"

type User struct {
	ID           int       `json:"id"`
	UUID         string    `json:"uuid"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Password     string    `json:"password"`
	AuthProvider string    `json:"auth_provider"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
