package client

// EndpointAssignmentRule is the endpoint to manage assignment rule (sysrule_assignment) records.
const EndpointAssignmentRule = "sysrule_assignment.do"

// AssignmentRule is the json response for an assignment rule in ServiceNow.
type AssignmentRule struct {
	BaseResult
	Name            string `json:"name"`
	Table           string `json:"table"`
	Active          bool   `json:"active,string"`
	Order           int    `json:"order,string"`
	Description     string `json:"description"`
	Condition       string `json:"condition"`
	AssignmentGroup string `json:"group"`
	User            string `json:"user"`
	Script          string `json:"script"`
	MatchFor        string `json:"match_conditions"`
}
