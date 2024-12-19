package models

import "time"

type Nasreload struct {
	NASIPAddress string    `json:"nas_ip_address"`
	ReloadTime   time.Time `json:"reload_time"`
}
