package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const serviceCatalogTitle = "title"
const serviceCatalogManager = "manager"
const serviceCatalogEditors = "editors"
const serviceCatalogDescription = "description"
const serviceCatalogBackgroundColor = "background_color"
const serviceCatalogDesktopImage = "desktop_image"
const serviceCatalogDesktopHomePage = "desktop_home_page"
const serviceCatalogDesktopContinueShopping = "desktop_continue_shopping"
const serviceCatalogActive = "active"
const serviceCatalogEnableWishList = "enable_wish_list"

// ResourceServieCatalog manages a Service Catalog object in ServiceNow.
func ResourceServiceCatalog() *schema.Resource {
	return &schema.Resource{
		Create: createResourceServiceCatalog,
		Read:   readResourceServiceCatalog,
		Update: updateResourceServiceCatalog,
		Delete: deleteResourceServiceCatalog,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			serviceCatalogTitle: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Title of the service catalog.",
			},
			serviceCatalogManager: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of manager user capable of performing updates, edits, and deletions again items and categories within service catalog.",
			},
			serviceCatalogEditors: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-seperated list of sys ids of editor users capable of editing and updating catalog categories and items.",
			},
			serviceCatalogDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-seperated list of sys ids of editor users capable of editing and updating catalog categories and items.",
			},
			serviceCatalogBackgroundColor: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "#FFFFFF",
				Description: "Background color of service catalog in hexidecimal format.",
			},
			serviceCatalogDesktopImage: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Path to desktop image displayed in service catalog",
			},
			serviceCatalogDesktopHomePage: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A custom to redirect to when user clicks the catalog breadcrumb.",
			},
			serviceCatalogDesktopContinueShopping: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A page name redirect for when a user clicks Continue Shopping.",
			},
			serviceCatalogActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', the service catalog will be visible to users.",
			},
			serviceCatalogEnableWishList: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If seto to 'true', wish list functionality will be visible in the service catalog.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceServiceCatalog(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalog := &client.ServiceCatalog{}
	if err := snowClient.GetObject(client.EndpointServiceCatalog, data.Id(), serviceCatalog); err != nil {
		data.SetId("")
		return err
	}

	resourceFromServiceCatalog(data, serviceCatalog)

	return nil
}

func createResourceServiceCatalog(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalog := resourceToServiceCatalog(data)
	if err := snowClient.CreateObject(client.EndpointServiceCatalog, serviceCatalog); err != nil {
		return err
	}

	resourceFromServiceCatalog(data, serviceCatalog)

	return readResourceServiceCatalog(data, serviceNowClient)
}

func updateResourceServiceCatalog(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointServiceCatalog, resourceToServiceCatalog(data)); err != nil {
		return err
	}

	return readResourceServiceCatalog(data, serviceNowClient)
}

func deleteResourceServiceCatalog(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return snowClient.DeleteObject(client.EndpointServiceCatalog, data.Id())
}

func resourceFromServiceCatalog(data *schema.ResourceData, serviceCatalog *client.ServiceCatalog) {
	data.SetId(serviceCatalog.ID)
	data.Set(serviceCatalogTitle, serviceCatalog.Title)
	data.Set(serviceCatalogManager, serviceCatalog.Manager)
	data.Set(serviceCatalogEditors, serviceCatalog.Editors)
	data.Set(serviceCatalogDescription, serviceCatalog.Description)
	data.Set(serviceCatalogBackgroundColor, serviceCatalog.BackgroundColor)
	data.Set(serviceCatalogDesktopImage, serviceCatalog.DesktopImage)
	data.Set(serviceCatalogDesktopHomePage, serviceCatalog.DesktopHomePage)
	data.Set(serviceCatalogDesktopContinueShopping, serviceCatalog.DesktopContinueShopping)
	data.Set(serviceCatalogActive, serviceCatalog.Active)
	data.Set(serviceCatalogEnableWishList, serviceCatalog.EnableWishList)
	data.Set(commonScope, serviceCatalog.Scope)
}

func resourceToServiceCatalog(data *schema.ResourceData) *client.ServiceCatalog {
	serviceCatalog := client.ServiceCatalog{
		Title:                   data.Get(serviceCatalogTitle).(string),
		Manager:                 data.Get(serviceCatalogManager).(string),
		Editors:                 data.Get(serviceCatalogEditors).(string),
		Description:             data.Get(serviceCatalogDescription).(string),
		BackgroundColor:         data.Get(serviceCatalogBackgroundColor).(string),
		DesktopImage:            data.Get(serviceCatalogDesktopImage).(string),
		DesktopHomePage:         data.Get(serviceCatalogDesktopHomePage).(string),
		DesktopContinueShopping: data.Get(serviceCatalogDesktopContinueShopping).(string),
		Active:                  data.Get(serviceCatalogActive).(bool),
		EnableWishList:          data.Get(serviceCatalogEnableWishList).(bool),
	}
	serviceCatalog.ID = data.Id()
	serviceCatalog.Scope = data.Get(commonScope).(string)
	return &serviceCatalog
}
