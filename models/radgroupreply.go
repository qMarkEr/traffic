package models

type Radgroupreply struct {
	Id        int64  `json:"id"`
	GroupName string `json:"group_name"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}
