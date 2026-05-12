package client

// EndpointBusinessRule is the endpoint to manage business rule records.
const EndpointBusinessRule = "sys_script.do"

// BusinessRule is the json response for a Business Rule in ServiceNow.
type BusinessRule struct {
	BaseResult
	Name            string `json:"name"`
	Table           string `json:"collection"`
	When            string `json:"when"`
	Order           int    `json:"order,string"`
	Active          bool   `json:"active,string"`
	Condition       string `json:"condition"`
	FilterCondition string `json:"filter_condition"`
	Script          string `json:"script"`
	Description     string `json:"description"`
	Advanced        bool   `json:"advanced,string"`
	ActionInsert    bool   `json:"action_insert,string"`
	ActionUpdate    bool   `json:"action_update,string"`
	ActionDelete    bool   `json:"action_delete,string"`
	ActionQuery     bool   `json:"action_query,string"`
	RoleConditions  string `json:"role_conditions"`
	Priority        int    `json:"priority,string"`
	AbortAction     bool   `json:"abort_action,string"`
	AddMessage      bool   `json:"add_message,string"`
	IsRest          bool   `json:"is_rest,string"`
	ClientCallable  bool   `json:"client_callable,string"`
}
