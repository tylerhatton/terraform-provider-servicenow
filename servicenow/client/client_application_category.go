package client

// EndpointApplicationCategory is the endpoint to manage application category records.
const EndpointApplicationCategory = "sys_app_category.do"

// ApplicationCategory represents the json response for a Application Category in ServiceNow.
type ApplicationCategory struct {
	BaseResult
	Name  string `json:"name"`
	Order int `json:"default_order,string"`
	Style string `json:"style"`
}
