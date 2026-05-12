package client

// EndpointUIScript is the endpoint to manage UI Script records.
const EndpointUIScript = "sys_ui_script.do"

// UIScript is the json response for a UI Script in ServiceNow.
type UIScript struct {
	BaseResult
	// Name maps to ServiceNow's "script_name" column, which is the user-visible
	// display name shown in the UI. APIName maps to ServiceNow's "name" column,
	// which is the system API name (typically computed/read-only). The mapping is
	// intentionally counterintuitive relative to the field identifiers.
	Name        string `json:"script_name"`
	APIName     string `json:"name,omitempty"`
	Description string `json:"description"`
	Script      string `json:"script"`
	Active      bool   `json:"active,string"`
	UIType      string `json:"ui_type"` // All: 10, Mobile: 1, Desktop 0
}
