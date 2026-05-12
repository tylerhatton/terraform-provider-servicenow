package client

// EndpointClientScript is the endpoint to manage client script records.
const EndpointClientScript = "sys_script_client.do"

// ClientScript is the json response for a Client Script in ServiceNow.
type ClientScript struct {
	BaseResult
	Name            string `json:"name"`
	Table           string `json:"table"`
	Type            string `json:"type"`
	FieldName       string `json:"field"`
	Script          string `json:"script"`
	Active          bool   `json:"active,string"`
	Description     string `json:"description"`
	Global          bool   `json:"global,string"`
	AppliesExtended bool   `json:"applies_extended,string"`
	Condition       string `json:"condition"`
	Messages        string `json:"messages"`
	View            string `json:"view"`
	UIType          string `json:"ui_type"`
}
