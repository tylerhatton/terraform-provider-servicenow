package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const serviceCatalogCategoryTitle = "title"
const serviceCatalogCategoryCatalog = "catalog"
const serviceCatalogCategoryLocation = "location"
const serviceCatalogCategoryDescription = "description"
const serviceCatalogCategoryDesktopImage = "desktop_image"
const serviceCatalogCategoryIcon = "icon"
const serviceCatalogCategoryHeaderIcon = "header_icon"
const serviceCatalogCategoryParent = "parent"
const serviceCatalogCategoryActive = "active"

// ResourceServieCatalogCategory manages a Service Catalog Category objects in ServiceNow.
func ResourceServiceCatalogCategory() *schema.Resource {
	return &schema.Resource{
		Create: createResourceServiceCatalogCategory,
		Read:   readResourceServiceCatalogCategory,
		Update: updateResourceServiceCatalogCategory,
		Delete: deleteResourceServiceCatalogCategory,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			serviceCatalogTitle: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Title of the service catalog category.",
			},
			serviceCatalogCategoryCatalog: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of service catalog the category is assigned to.",
			},
			serviceCatalogCategoryLocation: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of location of service catalog category",
			},
			serviceCatalogCategoryDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of service catalog category",
			},
			serviceCatalogCategoryDesktopImage: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Path to image associated with service catalog category.",
			},
			serviceCatalogCategoryIcon: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Path to service catalog category icon.",
			},
			serviceCatalogCategoryHeaderIcon: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Path to service catalog category header icon.",
			},
			serviceCatalogCategoryParent: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of parent category.",
			},
			serviceCatalogCategoryActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', the service catalog category will be visible to users.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceServiceCatalogCategory(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogCategory := &client.ServiceCatalogCategory{}
	if err := snowClient.GetObject(client.EndpointServiceCatalogCategory, data.Id(), serviceCatalogCategory); err != nil {
		data.SetId("")
		return err
	}

	resourceFromServiceCatalogCategory(data, serviceCatalogCategory)

	return nil
}

func createResourceServiceCatalogCategory(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogCategory := resourceToServiceCatalogCategory(data)
	if err := snowClient.CreateObject(client.EndpointServiceCatalogCategory, serviceCatalogCategory); err != nil {
		return err
	}

	resourceFromServiceCatalogCategory(data, serviceCatalogCategory)

	return readResourceServiceCatalogCategory(data, serviceNowClient)
}

func updateResourceServiceCatalogCategory(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointServiceCatalogCategory, resourceToServiceCatalogCategory(data)); err != nil {
		return err
	}

	return readResourceServiceCatalogCategory(data, serviceNowClient)
}

func deleteResourceServiceCatalogCategory(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return snowClient.DeleteObject(client.EndpointServiceCatalogCategory, data.Id())
}

func resourceFromServiceCatalogCategory(data *schema.ResourceData, serviceCatalogCategory *client.ServiceCatalogCategory) {
	data.SetId(serviceCatalogCategory.ID)
	data.Set(serviceCatalogCategoryTitle, serviceCatalogCategory.Title)
	data.Set(serviceCatalogCategoryCatalog, serviceCatalogCategory.Catalog)
	data.Set(serviceCatalogCategoryLocation, serviceCatalogCategory.Location)
	data.Set(serviceCatalogCategoryDescription, serviceCatalogCategory.Description)
	data.Set(serviceCatalogCategoryDesktopImage, serviceCatalogCategory.DesktopImage)
	data.Set(serviceCatalogCategoryIcon, serviceCatalogCategory.Icon)
	data.Set(serviceCatalogCategoryHeaderIcon, serviceCatalogCategory.HeaderIcon)
	data.Set(serviceCatalogCategoryParent, serviceCatalogCategory.Parent)
	data.Set(serviceCatalogCategoryActive, serviceCatalogCategory.Active)
	data.Set(commonScope, serviceCatalogCategory.Scope)
}

func resourceToServiceCatalogCategory(data *schema.ResourceData) *client.ServiceCatalogCategory {
	serviceCatalogCategory := client.ServiceCatalogCategory{
		Title:        data.Get(serviceCatalogCategoryTitle).(string),
		Catalog:      data.Get(serviceCatalogCategoryCatalog).(string),
		Location:     data.Get(serviceCatalogCategoryLocation).(string),
		Description:  data.Get(serviceCatalogCategoryDescription).(string),
		DesktopImage: data.Get(serviceCatalogCategoryDesktopImage).(string),
		Icon:         data.Get(serviceCatalogCategoryIcon).(string),
		HeaderIcon:   data.Get(serviceCatalogCategoryHeaderIcon).(string),
		Parent:       data.Get(serviceCatalogCategoryParent).(string),
		Active:       data.Get(serviceCatalogCategoryActive).(bool),
	}
	serviceCatalogCategory.ID = data.Id()
	serviceCatalogCategory.Scope = data.Get(commonScope).(string)
	return &serviceCatalogCategory
}
