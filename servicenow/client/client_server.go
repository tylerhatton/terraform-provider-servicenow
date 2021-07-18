package client

// EndpointServer is the endpoint to manage server entried in the CDMB.
const EndpointServer = "cmdb_ci_server.do"

// Server is the json response for a cmdb server entry in ServiceNow.
type Server struct {
	BaseResult
	Name            string `json:"name"`
	Company         string `json:"company"`
	AssetTag        string `json:"asset_tag"`
	SerialNumber    string `json:"serial_number"`
	Manufacturer    string `json:"manufacturer"`
	ModelId         string `json:"model_id"`
	AssignedTo      string `json:"assigned_to"`
	OsDomain        string `json:"os_domain"`
	Ram             string `json:"ram"`
	OperatingSystem string `json:"os"`
	CpuManufacturer string `json:"cpu_manufacturer"`
	OsVersion       string `json:"os_version"`
	CpuType         string `json:"cpu_type"`
	OsServicePack   string `json:"os_service_pack"`
	CpuSpeed        string `json:"cpu_speed"`
	DnsDomain       string `json:"dns_domain"`
	CpuCount        string `json:"cpu_count"`
	DiskSpace       string `json:"disk_space"`
	CpuCoreCount    string `json:"cpu_core_count"`
	Description     string `json:"short_description"`
	IpAddress       string `json:"ip_address"`
}
