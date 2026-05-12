package client

// EndpointDataLookup is the endpoint to manage data lookup (dl_definition) records.
const EndpointDataLookup = "dl_definition.do"

// DataLookup is the json response for a data lookup definition in ServiceNow.
//
// Note: dl_definition has a narrow column set: name, source_table,
// matcher_table, run_on_*. The legacy condition/matcher_fields/setter_fields/
// description/order fields are not part of the table and are intentionally
// omitted.
type DataLookup struct {
	BaseResult
	Name            string `json:"name"`
	Table           string `json:"source_table"`
	LookupTable     string `json:"matcher_table"`
	Active          bool   `json:"active,string"`
	RunOnInsert     bool   `json:"run_on_insert,string"`
	RunOnUpdate     bool   `json:"run_on_update,string"`
	RunOnFormChange bool   `json:"run_on_form_change,string"`
}
