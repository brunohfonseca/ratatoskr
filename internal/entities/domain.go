package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	Domain    string             `bson:"domain,omitempty" json:"domain,omitempty"`
	Port      int                `bson:"port,omitempty" json:"port,omitempty"`
	Enabled   bool               `bson:"enabled,omitempty" json:"enabled,omitempty"`
	LastCheck time.Time          `bson:"last_check,omitempty" json:"last_check,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
