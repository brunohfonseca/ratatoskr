package models

type AlertChannel struct {
	ID      int                    `json:"id"`
	Type    string                 `json:"type" index:""`
	Name    string                 `json:"name" index:"unique"`
	Config  map[string]interface{} `json:"config"`
	Enabled bool                   `json:"enabled"`
}

type AlertGroup struct {
	ID         int    `json:"id"`
	Name       string `json:"name" index:"unique"`
	ChannelIDs []int  `json:"channel_ids"`
	Enabled    bool   `json:"enabled"`
}
