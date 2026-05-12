package client

// EndpointMidServer is the endpoint to manage MID server records.
const EndpointMidServer = "ecc_agent.do"

// MidServer represents a MID server (ECC agent) record in ServiceNow. Note: MID server
// records are normally created by installing the MID server agent which then registers
// itself with ServiceNow. Creating these via the JSONv2 API is supported but typically
// only useful for placeholder/reference records.
type MidServer struct {
	BaseResult
	Name          string `json:"name"`
	HostName      string `json:"host_name"`
	Status        string `json:"status"`
	Version       string `json:"version"`
	Validated     bool   `json:"validated,string"`
	OSName        string `json:"host_os_distribution"`
	OSVersion     string `json:"host_os_version"`
	Description   string `json:"description,omitempty"`
	MidUser       string `json:"user_name"`
	Started       string `json:"started"`
	LastRefresh   string `json:"last_refreshed"`
	AgentType     string `json:"type"`
	LinuxUserName string `json:"linux_user_name,omitempty"`
}
