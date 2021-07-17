package client

// EndpointServiceCatalogVariable is the endpoint to manage service catalog variables.
const EndpointServiceCatalogVariable = "item_option_new.do"

// ServiceCatalogVariable is the json response for a service catalog variable in ServiceNow.
type ServiceCatalogVariable struct {
	BaseResult
	Name               string `json:"name"`
	Question           string `json:"question_text"`
	Tooltip            string `json:"tooltip"`
	HelpTag            string `json:"help_tag"`
	HelpText           string `json:"help_text"`
	Instructions       string `json:"instructions"`
	DefaultValue       string `json:"default_value"`
	Type               string `json:"type"`
	CatalogItem        string `json:"cat_item"`
	Order              string `json:"order"`
	ListTable          string `json:"list_table"`
	LookupTable        string `json:"lookup_table"`
	LookupValue        string `json:"lookup_value"`
	Reference          string `json:"reference"`
	ReferenceQualifier string `json:"reference_qual"`
	ShowHelp           bool   `json:"show_help,string"`
	Mandatory          bool   `json:"mandatory,string"`
	ReadOnly           bool   `json:"read_only,string"`
	Hidden             bool   `json:"hidden,string"`
	Active             bool   `json:"active,string"`
}
