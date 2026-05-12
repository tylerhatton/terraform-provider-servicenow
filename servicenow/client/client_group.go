package client

// EndpointGroup is the endpoint to manage user group records.
const EndpointGroup = "sys_user_group.do"

// Group is the json response for a user group in ServiceNow.
type Group struct {
	BaseResult
	Name            string `json:"name"`
	Description     string `json:"description"`
	Email           string `json:"email"`
	Manager         string `json:"manager"`
	Parent          string `json:"parent"`
	CostCenter      string `json:"cost_center"`
	Type            string `json:"type"`
	Active          bool   `json:"active,string"`
	DefaultAssignee string `json:"default_assignee"`
	IncludeMembers  bool   `json:"include_members,string"`
	ExcludeManager  bool   `json:"exclude_manager,string"`
}
