package models

import "time"

type Radacct struct {
	RadAcctId           int64     `json:"rad_acct_id"`
	AcctSessionId       string    `json:"acct_session_id"`
	AcctUniqueId        string    `json:"acct_unique_id"`
	UserName            string    `json:"user_name"`
	Realm               string    `json:"realm"`
	NASIPAddress        string    `json:"nas_ip_address"`
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
	FramedIPAddress     string    `json:"framed_ip_address"`
	FramedIPv6Address   string    `json:"framed_ipv6_address"`
	FramedIPv6Prefix    string    `json:"framed_ipv6_prefix"`
	FramedInterfaceId   string    `json:"framed_interface_id"`
	DelegatedIPv6Prefix string    `json:"delegated_ipv6_prefix"`
	Class               string    `json:"class"`
}
