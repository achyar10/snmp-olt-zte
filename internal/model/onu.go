package model

type OltConfig struct {
	BaseOID                   string
	OnuIDNameOID              string
	OnuTypeOID                string
	OnuSerialNumberOID        string
	OnuRxPowerOID             string
	OnuTxPowerOID             string
	OnuStatusOID              string
	OnuIPAddressOID           string
	OnuDescriptionOID         string
	OnuLastOnlineOID          string
	OnuLastOfflineOID         string
	OnuLastOfflineReasonOID   string
	OnuGponOpticalDistanceOID string
}

type ONUInfo struct {
	ID   string `json:"onu_id"`
	Name string `json:"name"`
}

type ONUInfoPerBoard struct {
	Board        int    `json:"board"`
	PON          int    `json:"pon"`
	ID           int    `json:"onu_id"`
	Name         string `json:"name"`
	OnuType      string `json:"onu_type"`
	SerialNumber string `json:"serial_number"`
	RXPower      string `json:"rx_power"`
	Status       string `json:"status"`
}

type ONUCustomerInfo struct {
	Board                int    `json:"board"`
	PON                  int    `json:"pon"`
	ID                   int    `json:"onu_id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	OnuType              string `json:"onu_type"`
	SerialNumber         string `json:"serial_number"`
	RXPower              string `json:"rx_power"`
	TXPower              string `json:"tx_power"`
	Status               string `json:"status"`
	IPAddress            string `json:"ip_address"`
	LastOnline           string `json:"last_online"`
	LastOffline          string `json:"last_offline"`
	Uptime               string `json:"uptime"`
	LastDownTimeDuration string `json:"last_down_time_duration"`
	LastOfflineReason    string `json:"offline_reason"`
	GponOpticalDistance  string `json:"gpon_optical_distance"`
}

type OnuID struct {
	Board int `json:"board"`
	PON   int `json:"pon"`
	ID    int `json:"onu_id"`
}

type OnuOnlyID struct {
	ID int `json:"onu_id"`
}

type SNMPWalkTask struct {
	BaseOID   string
	TargetOID string
	BoardID   int
	PON       int
}

type OnuSerialNumber struct {
	Board        int    `json:"board"`
	PON          int    `json:"pon"`
	ID           int    `json:"onu_id"`
	SerialNumber string `json:"serial_number"`
}

type PaginationResult struct {
	OnuInformationList []ONUInfoPerBoard
	Count              int
}

type TelnetRequest struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ONUItem struct {
	OltIndex     string `json:"olt_index"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Status       string `json:"status"`
}

type ONUStatus struct {
	ID         int    `json:"id"`
	Status     string `json:"status"`
	AdminState string `json:"admin_state,omitempty"`
	OMCCState  string `json:"omcc_state,omitempty"`
	PhaseState string `json:"phase_state,omitempty"`
	Channel    string `json:"channel,omitempty"`
}

type ActivateONURequest struct {
	OLTIndex     string `json:"olt_index"`
	SerialNumber string `json:"serial_number"`
	Region       string `json:"region"`
	Code         string `json:"code"`
	VlanID       int    `json:"vlan_id,omitempty"`
	Onu          *int   `json:"onu,omitempty"` // pointer to know if it’s provided or not
}
