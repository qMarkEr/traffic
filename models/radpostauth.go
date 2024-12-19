package models

import "time"

type Radpostauth struct {
	Id               int64     `json:"id"`
	Username         string    `json:"username"`
	Pass             string    `json:"pass"`
	Reply            string    `json:"reply"`
	CalledStationId  string    `json:"called_station_id"`
	CallingStationId string    `json:"calling_station_id"`
	Authdate         time.Time `json:"authdate"`
	Class            string    `json:"class"`
}
