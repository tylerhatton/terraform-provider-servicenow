package client

// EndpointDictionary is the endpoint to manage dictionary entries (column definitions) in ServiceNow.
const EndpointDictionary = "sys_dictionary.do"

// Dictionary is the json response for a dictionary entry (table column definition) in ServiceNow.
type Dictionary struct {
	BaseResult
	Name             string `json:"name"`
	Element          string `json:"element"`
	ColumnLabel      string `json:"column_label,omitempty"`
	InternalType     string `json:"internal_type,omitempty"`
	MaxLength        string `json:"max_length,omitempty"`
	Mandatory        bool   `json:"mandatory,string"`
	ReadOnly         bool   `json:"read_only,string"`
	Active           bool   `json:"active,string"`
	Display          bool   `json:"display,string"`
	Unique           bool   `json:"unique,string"`
	DefaultValue     string `json:"default_value,omitempty"`
	Comments         string `json:"comments,omitempty"`
	Reference        string `json:"reference,omitempty"`
	DynamicCreation  bool   `json:"dynamic_creation,string"`
	Dependent        string `json:"dependent,omitempty"`
	DependentOnField string `json:"dependent_on_field,omitempty"`
	Choice           string `json:"choice,omitempty"`
}
