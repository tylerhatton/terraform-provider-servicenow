package client

// EndpointAlias is the endpoint to manage connection and credential alias records.
const EndpointAlias = "sys_alias.do"

// Alias is the json response for a connection and credential in ServiceNow.
type Alias struct {
	BaseResult
	Name                  string `json:"name"`
	ParentAlias           string `json:"parent"`
	Type                  string `json:"type"`
	ConnectionType        string `json:"connection_type"`
	MultipleActions       bool   `json:"multiple_actions,string"`
	RetryPolicy           string `json:"retry_policy"`
	ConfigurationTemplate string `json:"configuration_template"`
}
