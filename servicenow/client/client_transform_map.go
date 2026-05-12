package client

// EndpointTransformMap is the endpoint to manage transform map (sys_transform_map) records.
const EndpointTransformMap = "sys_transform_map.do"

// TransformMap is the json response for an import transform map in ServiceNow.
//
// Note: sys_transform_map does not include a `description` column.
type TransformMap struct {
	BaseResult
	Name                   string `json:"name"`
	SourceTable            string `json:"source_table"`
	TargetTable            string `json:"target_table"`
	Active                 bool   `json:"active,string"`
	RunBusinessRules       bool   `json:"run_business_rules,string"`
	EnforceMandatoryFields string `json:"enforce_mandatory_fields"`
	CopyEmptyFields        bool   `json:"copy_empty_fields,string"`
	Order                  int    `json:"order,string"`
	Script                 string `json:"script"`
}
