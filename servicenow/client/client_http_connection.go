package client

// EndpointHttpConnection is the endpoint to manage system property categories records.
const EndpointHttpConnection = "http_connection.do"

// HttpConnection is the json response for a system property in ServiceNow.
type HttpConnection struct {
	BaseResult
	Name            string `json:"name"`
	Active          bool   `json:"active,string"`
	Credential      string `json:"credential"`
	ConnectionAlias string `json:"connection_alias"`
	ConnectionUrl   string `json:"connection_url"`
	UseMidServer    bool   `json:"use_mid,string"`
	MidSelection    string `json:"mid_selection"`
	MidServer       string `json:"mid_server"`
}
