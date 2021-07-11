package client

// EndpointServiceCatalog is the endpoint to manage service catalog records.
const EndpointServiceCatalog = "sc_catalog.do"

// ServiceCatalog is the json response for a service catalog in ServiceNow.
type ServiceCatalog struct {
	BaseResult
	Title                   string `json:"title"`
	Manager                 string `json:"manager"`
	Editors                 string `json:"editors"`
	Description             string `json:"description"`
	BackgroundColor         string `json:"background_color"`
	DesktopImage            string `json:"desktop_image"`
	DesktopHomePage         string `json:"desktop_home_page"`
	DesktopContinueShopping string `json:"desktop_continue_shopping"`
	Active                  bool   `json:"active,string"`
	EnableWishList          bool   `json:"enable_wish_list,string"`
}
