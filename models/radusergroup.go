package models

type Radusergroup struct {
	Id        int64  `json:"id"`
	UserName  string `json:"user_name"`
	GroupName string `json:"group_name"`
	Priority  int    `json:"priority"`
}
