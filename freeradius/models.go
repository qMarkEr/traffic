// models.go
package freeradius

import (
	"net"
	"time"
)

// Структура для таблицы 'radacct'
type RadAcct struct {
	RadAcctId           int64     `json:"radacct_id"`
	AcctSessionId       string    `json:"acct_session_id"`
	AcctUniqueId        string    `json:"acct_unique_id"`
	UserName            string    `json:"username"`
	Realm               string    `json:"realm"`
	NASIPAddress        net.IP    `json:"nas_ip_address"`
	NASPortId           string    `json:"nas_port_id"`
	NASPortType         string    `json:"nas_port_type"`
	AcctStartTime       time.Time `json:"acct_start_time"`
	AcctUpdateTime      time.Time `json:"acct_update_time"`
	AcctStopTime        time.Time `json:"acct_stop_time"`
	AcctInterval        int64     `json:"acct_interval"`
	AcctSessionTime     int64     `json:"acct_session_time"`
	AcctAuthentic       string    `json:"acct_authentic"`
	ConnectInfoStart    string    `json:"connect_info_start"`
	ConnectInfoStop     string    `json:"connect_info_stop"`
	AcctInputOctets     int64     `json:"acct_input_octets"`
	AcctOutputOctets    int64     `json:"acct_output_octets"`
	CalledStationId     string    `json:"called_station_id"`
	CallingStationId    string    `json:"calling_station_id"`
	AcctTerminateCause  string    `json:"acct_terminate_cause"`
	ServiceType         string    `json:"service_type"`
	FramedProtocol      string    `json:"framed_protocol"`
	FramedIPAddress     net.IP    `json:"framed_ip_address"`
	FramedIPv6Address   net.IP    `json:"framed_ipv6_address"`
	FramedIPv6Prefix    net.IP    `json:"framed_ipv6_prefix"`
	FramedInterfaceId   string    `json:"framed_interface_id"`
	DelegatedIPv6Prefix net.IP    `json:"delegated_ipv6_prefix"`
	Class               string    `json:"class"`
}

// Структура для таблицы 'radcheck'
type RadCheck struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}

// Структура для таблицы 'radgroupcheck'
type RadGroupCheck struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}

// Структура для таблицы 'radgroupreply'
type RadGroupReply struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}

// Структура для таблицы 'radreply'
type RadReply struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}

// Структура для таблицы 'radusergroup'
type RadUserGroup struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	GroupName string `json:"group_name"`
	Priority  int    `json:"priority"`
}

// Структура для таблицы 'radpostauth'
type RadPostAuth struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	Pass             string    `json:"pass"`
	Reply            string    `json:"reply"`
	CalledStationId  string    `json:"called_station_id"`
	CallingStationId string    `json:"calling_station_id"`
	AuthDate         time.Time `json:"auth_date"`
	Class            string    `json:"class"`
}

// Структура для таблицы 'nas'
type Nas struct {
	ID          int    `json:"id"`
	NasName     string `json:"nas_name"`
	ShortName   string `json:"short_name"`
	Type        string `json:"type"`
	Ports       int    `json:"ports"`
	Secret      string `json:"secret"`
	Server      string `json:"server"`
	Community   string `json:"community"`
	Description string `json:"description"`
}

// Структура для таблицы 'nasreload'
type NasReload struct {
	NASIPAddress net.IP    `json:"nas_ip_address"`
	ReloadTime   time.Time `json:"reload_time"`
}
