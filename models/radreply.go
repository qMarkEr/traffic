package models

type Radreply struct {
	Id        int64  `json:"id"`
	UserName  string `json:"user_name"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}
