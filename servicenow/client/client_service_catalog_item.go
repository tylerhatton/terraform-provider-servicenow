package client

// EndpointServiceCatalogItem is the endpoint to manage service catalog item records.
const EndpointServiceCatalogItem = "sc_cat_item.do"

// ServiceCatalogItem is the json response for a service catalog item in ServiceNow.
type ServiceCatalogItem struct {
	BaseResult
	Name                string `json:"name"`
	ServiceCatalogs     string `json:"sc_catalogs"`
	Category            string `json:"category"`
	ShortDescription    string `json:"short_description"`
	Description         string `json:"description"`
	HideAddToCart       bool   `json:"no_cart_v2,string"`
	HideQuantity        bool   `json:"no_quantity_v2,string"`
	HideDeliveryTime    bool   `json:"no_delivery_time_v2,string"`
	HideAddToWishlist   bool   `json:"no_wishlist_v2,string"`
	HideAttachment      bool   `json:"no_attachment_v2,string"`
	MandatoryAttachment bool   `json:"mandatory_attachment,string"`
	Active              bool   `json:"active,string"`
}
