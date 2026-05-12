package client

// EndpointGroupMember is the endpoint to manage group membership records.
const EndpointGroupMember = "sys_user_grmember.do"

// GroupMember is the json response for a user-group membership in ServiceNow.
type GroupMember struct {
	BaseResult
	User  string `json:"user"`
	Group string `json:"group"`
}
