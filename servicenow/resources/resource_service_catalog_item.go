package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const serviceCatalogItemName = "name"
const serviceCatalogItemServiceCatalogs = "service_catalogs"
const serviceCatalogItemCategory = "category"
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
		Description: "`servicenow_service_catalog_item` manages a service catalog item configuration within ServiceNow.",

		CreateContext: createResourceServiceCatalogItem,
		ReadContext:   readResourceServiceCatalogItem,
		UpdateContext: updateResourceServiceCatalogItem,
		DeleteContext: deleteResourceServiceCatalogItem,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			serviceCatalogItemName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of service catalog item.",
			},
			serviceCatalogItemServiceCatalogs: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-seperated list of service catalogs the service catalog item is assigned to.",
			},
			serviceCatalogItemCategory: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Service catalog category the service catalog item is assigned to.",
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

func readResourceServiceCatalogItem(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogItem := &client.ServiceCatalogItem{}
	if err := snowClient.GetObject(client.EndpointServiceCatalogItem, data.Id(), serviceCatalogItem); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromServiceCatalogItem(data, serviceCatalogItem)

	return nil
}

func createResourceServiceCatalogItem(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogItem := resourceToServiceCatalogItem(data)
	if err := snowClient.CreateObject(client.EndpointServiceCatalogItem, serviceCatalogItem); err != nil {
		return diag.FromErr(err)
	}

	resourceFromServiceCatalogItem(data, serviceCatalogItem)

	return readResourceServiceCatalogItem(ctx, data, serviceNowClient)
}

func updateResourceServiceCatalogItem(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointServiceCatalogItem, resourceToServiceCatalogItem(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceServiceCatalogItem(ctx, data, serviceNowClient)
}

func deleteResourceServiceCatalogItem(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointServiceCatalogItem, data.Id()))
}

func resourceFromServiceCatalogItem(data *schema.ResourceData, serviceCatalogItem *client.ServiceCatalogItem) {
	data.SetId(serviceCatalogItem.ID)
	data.Set(serviceCatalogItemName, serviceCatalogItem.Name)
	data.Set(serviceCatalogItemServiceCatalogs, serviceCatalogItem.ServiceCatalogs)
	data.Set(serviceCatalogItemCategory, serviceCatalogItem.Category)
	data.Set(serviceCatalogItemShortDescription, serviceCatalogItem.ShortDescription)
	data.Set(serviceCatalogItemDescription, serviceCatalogItem.Description)
	data.Set(serviceCatalogItemHideAddToCart, serviceCatalogItem.HideAddToCart)
	data.Set(serviceCatalogItemHideQuantity, serviceCatalogItem.HideQuantity)
	data.Set(serviceCatalogItemHideDeliveryTime, serviceCatalogItem.HideDeliveryTime)
	data.Set(serviceCatalogItemHideAddToWishlist, serviceCatalogItem.HideAddToWishlist)
	data.Set(serviceCatalogItemHideAttachment, serviceCatalogItem.HideAttachment)
	data.Set(serviceCatalogItemMandatoryAttachment, serviceCatalogItem.MandatoryAttachment)
	data.Set(serviceCatalogItemActive, serviceCatalogItem.Active)
	data.Set(commonScope, serviceCatalogItem.Scope)
}

func resourceToServiceCatalogItem(data *schema.ResourceData) *client.ServiceCatalogItem {
	serviceCatalogItem := client.ServiceCatalogItem{
		Name:                data.Get(serviceCatalogItemName).(string),
		ServiceCatalogs:     data.Get(serviceCatalogItemServiceCatalogs).(string),
		Category:            data.Get(serviceCatalogItemCategory).(string),
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
