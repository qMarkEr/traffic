package models

type Nas struct {
	Id          int64  `json:"id"`
	Nasname     string `json:"nasname"`
	Shortname   string `json:"shortname"`
	Type        string `json:"type"`
	Ports       int    `json:"ports"`
	Secret      string `json:"secret"`
	Server      string `json:"server"`
	Community   string `json:"community"`
	Description string `json:"description"`
}
