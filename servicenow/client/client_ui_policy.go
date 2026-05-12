package client

// EndpointUIPolicy is the endpoint to manage UI policy records.
const EndpointUIPolicy = "sys_ui_policy.do"

// UIPolicy is the json response for a UI Policy in ServiceNow.
type UIPolicy struct {
	BaseResult
	ShortDescription string `json:"short_description"`
	Table            string `json:"table"`
	Active           bool   `json:"active,string"`
	Conditions       string `json:"conditions"`
	ReverseIfFalse   bool   `json:"reverse_if_false,string"`
	RunScripts       bool   `json:"run_scripts,string"`
	ScriptTrue       string `json:"script_true"`
	ScriptFalse      string `json:"script_false"`
	OnLoad           bool   `json:"on_load,string"`
	Global           bool   `json:"global,string"`
	Order            int    `json:"order,string"`
	View             string `json:"view"`
	Description      string `json:"description"`
}
