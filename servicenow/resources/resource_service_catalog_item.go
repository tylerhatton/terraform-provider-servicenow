package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const serviceCatalogItemName = "name"
const serviceCatalogItemShortDescription = "short_description"
const serviceCatalogItemDescription = "description"
const serviceCatalogItemHideAddToCart = "no_cart"
const serviceCatalogItemHideQuantity = "no_quantity"
const serviceCatalogItemHideDeliveryTime = "no_delivery_time"
const serviceCatalogItemHideAddToWishlist = "no_wishlist"
const serviceCatalogItemHideAttachment = "no_attachment"
const serviceCatalogItemMandatoryAttachment = "mandatory_attachment"
const serviceCatalogItemActive = "active"

// ResourceServiceCatalogItem manages a Service Catalog Item in ServiceNow.
func ResourceServiceCatalogItem() *schema.Resource {
	return &schema.Resource{
		Create: createResourceServiceCatalogItem,
		Read:   readResourceServiceCatalogItem,
		Update: updateResourceServiceCatalogItem,
		Delete: deleteResourceServiceCatalogItem,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			serviceCatalogItemName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of service catalog item.",
			},
			serviceCatalogItemShortDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Short description of service catalog item.",
			},
			serviceCatalogItemDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of service catalog item.",
			},
			serviceCatalogItemHideAddToCart: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will hide the add to cart button from the service catalog item.",
			},
			serviceCatalogItemHideQuantity: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will hide the quantity parameter from the service catalog item.",
			},
			serviceCatalogItemHideDeliveryTime: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will hide the delivery time from the service catalog item.",
			},
			serviceCatalogItemHideAddToWishlist: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will hide the add to wishlist button from the service catalog item.",
			},
			serviceCatalogItemHideAttachment: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will hide the attachment button from the service catalog item.",
			},
			serviceCatalogItemMandatoryAttachment: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will require an attachment when a service catalog item is used.",
			},
			serviceCatalogItemActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', this property will make a service catalog visible for use.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceServiceCatalogItem(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogItem := &client.ServiceCatalogItem{}
	if err := snowClient.GetObject(client.EndpointServiceCatalogItem, data.Id(), serviceCatalogItem); err != nil {
		data.SetId("")
		return err
	}

	resourceFromServiceCatalogItem(data, serviceCatalogItem)

	return nil
}

func createResourceServiceCatalogItem(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogItem := resourceToServiceCatalogItem(data)
	if err := snowClient.CreateObject(client.EndpointServiceCatalogItem, serviceCatalogItem); err != nil {
		return err
	}

	resourceFromServiceCatalogItem(data, serviceCatalogItem)

	return readResourceServiceCatalogItem(data, serviceNowClient)
}

func updateResourceServiceCatalogItem(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointServiceCatalogItem, resourceToServiceCatalogItem(data)); err != nil {
		return err
	}

	return readResourceServiceCatalogItem(data, serviceNowClient)
}

func deleteResourceServiceCatalogItem(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return snowClient.DeleteObject(client.EndpointServiceCatalogItem, data.Id())
}

func resourceFromServiceCatalogItem(data *schema.ResourceData, serviceCatalogItem *client.ServiceCatalogItem) {
	data.SetId(serviceCatalogItem.ID)
	data.Set(serviceCatalogItemName, serviceCatalogItem.Name)
	data.Set(serviceCatalogItemShortDescription, serviceCatalogItem.ShortDescription)
	data.Set(serviceCatalogItemDescription, serviceCatalogItem.Description)
	data.Set(serviceCatalogItemHideAddToCart, serviceCatalogItem.HideAddToCart)
	data.Set(serviceCatalogItemHideQuantity, serviceCatalogItem.HideQuantity)
	data.Set(serviceCatalogItemHideDeliveryTime, serviceCatalogItem.HideDeliveryTime)
	data.Set(serviceCatalogItemHideAddToWishlist, serviceCatalogItem.HideAddToWishlist)
	data.Set(serviceCatalogItemHideAttachment, serviceCatalogItem.HideAttachment)
	data.Set(serviceCatalogItemMandatoryAttachment, serviceCatalogItem.MandatoryAttachment)
	data.Set(serviceCatalogItemActive, serviceCatalogItem.Active)
}

func resourceToServiceCatalogItem(data *schema.ResourceData) *client.ServiceCatalogItem {
	serviceCatalogItem := client.ServiceCatalogItem{
		Name:                data.Get(serviceCatalogItemName).(string),
		ShortDescription:    data.Get(serviceCatalogItemShortDescription).(string),
		Description:         data.Get(serviceCatalogItemDescription).(string),
		HideAddToCart:       data.Get(serviceCatalogItemHideAddToCart).(bool),
		HideQuantity:        data.Get(serviceCatalogItemHideQuantity).(bool),
		HideDeliveryTime:    data.Get(serviceCatalogItemHideDeliveryTime).(bool),
		HideAddToWishlist:   data.Get(serviceCatalogItemHideAddToWishlist).(bool),
		HideAttachment:      data.Get(serviceCatalogItemHideAttachment).(bool),
		MandatoryAttachment: data.Get(serviceCatalogItemMandatoryAttachment).(bool),
		Active:              data.Get(serviceCatalogItemActive).(bool),
	}
	serviceCatalogItem.ID = data.Id()
	serviceCatalogItem.Scope = data.Get(commonScope).(string)
	return &serviceCatalogItem
}
