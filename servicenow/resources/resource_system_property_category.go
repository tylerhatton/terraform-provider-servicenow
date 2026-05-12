package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const systemPropertyCategoryName = "name"
const systemPropertyCategoryTitleHTML = "title_html"

// ResourceSystemPropertyCategory manages a System Property Category in ServiceNow.
func ResourceSystemPropertyCategory() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_system_property_category` manages a system property category within ServiceNow.",

		CreateContext: createResourceSystemPropertyCategory,
		ReadContext:   readResourceSystemPropertyCategory,
		UpdateContext: updateResourceSystemPropertyCategory,
		DeleteContext: deleteResourceSystemPropertyCategory,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			systemPropertyCategoryName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the category.",
			},
			systemPropertyCategoryTitleHTML: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The HTML displayed at the top of the page when configuring properties for this category.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceSystemPropertyCategory(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	systemPropertyCategory := &client.SystemPropertyCategory{}
	if err := snowClient.GetObject(client.EndpointSystemPropertyCategory, data.Id(), systemPropertyCategory); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromSystemPropertyCategory(data, systemPropertyCategory)

	return nil
}

func createResourceSystemPropertyCategory(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	systemPropertyCategory := resourceToSystemPropertyCategory(data)
	if err := snowClient.CreateObject(client.EndpointSystemPropertyCategory, systemPropertyCategory); err != nil {
		return diag.FromErr(err)
	}

	resourceFromSystemPropertyCategory(data, systemPropertyCategory)

	return readResourceSystemPropertyCategory(ctx, data, serviceNowClient)
}

func updateResourceSystemPropertyCategory(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointSystemPropertyCategory, resourceToSystemPropertyCategory(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceSystemPropertyCategory(ctx, data, serviceNowClient)
}

func deleteResourceSystemPropertyCategory(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointSystemPropertyCategory, data.Id()))
}

func resourceFromSystemPropertyCategory(data *schema.ResourceData, systemPropertyCategory *client.SystemPropertyCategory) {
	data.SetId(systemPropertyCategory.ID)
	data.Set(systemPropertyCategoryName, systemPropertyCategory.Name)
	data.Set(systemPropertyCategoryTitleHTML, systemPropertyCategory.TitleHTML)
	data.Set(commonScope, systemPropertyCategory.Scope)
}

func resourceToSystemPropertyCategory(data *schema.ResourceData) *client.SystemPropertyCategory {
	systemPropertyCategory := client.SystemPropertyCategory{
		Name:      data.Get(systemPropertyCategoryName).(string),
		TitleHTML: data.Get(systemPropertyCategoryTitleHTML).(string),
	}
	systemPropertyCategory.ID = data.Id()
	systemPropertyCategory.Scope = data.Get(commonScope).(string)
	return &systemPropertyCategory
}
