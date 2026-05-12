package client

// EndpointUserRole is the endpoint to manage user-to-role associations.
const EndpointUserRole = "sys_user_has_role.do"

// UserRole is the json response for a user-role association in ServiceNow.
type UserRole struct {
	BaseResult
	User      string `json:"user"`
	Role      string `json:"role"`
	Inherited bool   `json:"inherited,string"`
	GrantedBy string `json:"granted_by"`
	State     string `json:"state"`
}
