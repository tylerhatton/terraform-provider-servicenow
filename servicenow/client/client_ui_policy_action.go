package client

// EndpointUIPolicyAction is the endpoint to manage UI policy action records.
const EndpointUIPolicyAction = "sys_ui_policy_action.do"

// UIPolicyAction is the json response for a UI Policy Action in ServiceNow.
type UIPolicyAction struct {
	BaseResult
	UIPolicy  string `json:"ui_policy"`
	FieldName string `json:"field"`
	Visible   string `json:"visible"`
	Mandatory string `json:"mandatory"`
	Disabled  string `json:"disabled"`
	Cleared   bool   `json:"cleared,string"`
}
