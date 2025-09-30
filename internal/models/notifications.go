package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AlertChannel struct {
	ID      primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Type    string                 `bson:"type" json:"type" index:""`
	Name    string                 `bson:"name" json:"name" index:"unique"`
	Config  map[string]interface{} `bson:"config" json:"config"`
	Enabled bool                   `bson:"enabled" json:"enabled"`
}

type AlertGroup struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name       string               `bson:"name" json:"name" index:"unique"`
	ChannelIDs []primitive.ObjectID `bson:"channel_ids" json:"channel_ids"`
	Enabled    bool                 `bson:"enabled" json:"enabled"`
}
