package client

// EndpointChoice is the endpoint to manage choice list values (sys_choice) in ServiceNow.
const EndpointChoice = "sys_choice.do"

// Choice is the json response for a choice list entry in ServiceNow.
type Choice struct {
	BaseResult
	Name           string `json:"name"`
	Element        string `json:"element"`
	Value          string `json:"value"`
	Label          string `json:"label,omitempty"`
	Sequence       string `json:"sequence,omitempty"`
	Hint           string `json:"hint,omitempty"`
	Inactive       bool   `json:"inactive,string"`
	DependentValue string `json:"dependent_value,omitempty"`
	Language       string `json:"language,omitempty"`
}
