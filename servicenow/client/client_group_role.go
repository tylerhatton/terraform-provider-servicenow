package client

// EndpointGroupRole is the endpoint to manage group-to-role associations.
const EndpointGroupRole = "sys_group_has_role.do"

// GroupRole is the json response for a group-role association in ServiceNow.
type GroupRole struct {
	BaseResult
	Group    string `json:"group"`
	Role     string `json:"role"`
	Inherits bool   `json:"inherits,string"`
}
