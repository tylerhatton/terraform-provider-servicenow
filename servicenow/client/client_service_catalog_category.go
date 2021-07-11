package client

// EndpointServiceCatalogCategory is the endpoint to manage service catalog categories records.
const EndpointServiceCatalogCategory = "sc_category.do"

// ServiceCatalogCategory is the json response for a service catalog category in ServiceNow.
type ServiceCatalogCategory struct {
	BaseResult
	Title        string `json:"title"`
	Catalog      string `json:"sc_catalog"`
	Location     string `json:"location"`
	Description  string `json:"description"`
	DesktopImage string `json:"image"`
	Icon         string `json:"icon"`
	HeaderIcon   string `json:"header_icon"`
	Parent       string `json:"parent"`
	Active       bool   `json:"active,string"`
}
