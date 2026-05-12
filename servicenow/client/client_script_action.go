package client

// EndpointScriptAction is the endpoint to manage script action (sysevent_script_action) records.
const EndpointScriptAction = "sysevent_script_action.do"

// ScriptAction is the json response for an event-driven script action in ServiceNow.
type ScriptAction struct {
	BaseResult
	Name        string `json:"name"`
	EventName   string `json:"event_name"`
	Script      string `json:"script"`
	Active      bool   `json:"active,string"`
	Description string `json:"description"`
	Order       int    `json:"order,string"`
	Condition   string `json:"condition_script"`
	Synchronous bool   `json:"synchronous,string"`
}
