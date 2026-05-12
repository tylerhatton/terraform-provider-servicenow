package client

// EndpointFlow is the endpoint to manage Flow Designer flow records.
const EndpointFlow = "sys_hub_flow.do"

// Flow represents a Flow Designer flow record in ServiceNow.
//
// NOTE: The full Flow Designer schema (steps, variables, triggers, snapshots,
// label cache, etc.) is intentionally not modelled here. Flow Designer flows are
// composed of many related tables (sys_hub_action_instance,
// sys_hub_trigger_instance, sys_hub_flow_input, sys_hub_flow_output, ...) and
// are normally authored interactively in the Flow Designer UI. This struct only
// covers the high level metadata exposed on the sys_hub_flow table itself, so
// you can manage flow shells / references via Terraform while still authoring
// the flow logic in the UI.
//
// The `Application` Terraform attribute maps to the BaseResult.Scope (sys_scope)
// column inherited from BaseResult.
type Flow struct {
	BaseResult
	Name         string `json:"name"`
	Description  string `json:"description"`
	Active       bool   `json:"active,string"`
	TriggerType  string `json:"trigger_type,omitempty"`
	Category     string `json:"type"`
	InternalName string `json:"internal_name"`
	Status       string `json:"status"`
}
