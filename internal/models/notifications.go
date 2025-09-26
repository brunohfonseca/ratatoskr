package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AlertChannel struct {
	Type    string                 // "telegram", "slack", "email"
	Config  map[string]interface{} // webhook_url, chat_id, etc
	Enabled bool
}

type AlertGroup struct {
	Name       string
	ChannelIDs []primitive.ObjectID
	Enabled    bool
}
