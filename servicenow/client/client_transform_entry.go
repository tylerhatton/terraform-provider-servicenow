package client

// EndpointTransformEntry is the endpoint to manage transform entry (sys_transform_entry) records.
const EndpointTransformEntry = "sys_transform_entry.do"

// TransformEntry is the json response for a transform map field entry in ServiceNow.
//
// Note: sys_transform_entry does not include `order` or `reference_table`
// columns. The reference value field is exposed as `reference_value_field`.
type TransformEntry struct {
	BaseResult
	Map                      string `json:"map"`
	SourceField              string `json:"source_field"`
	TargetField              string `json:"target_field"`
	Coalesce                 bool   `json:"coalesce,string"`
	UseSourceScript          bool   `json:"use_source_script,string"`
	SourceScript             string `json:"source_script"`
	ReferencedValueFieldName string `json:"reference_value_field"`
}
