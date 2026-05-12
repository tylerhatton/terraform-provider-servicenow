package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const widgetDependencyName = "name"
const widgetDependencyModule = "module"
const widgetDependencyPageLoad = "page_load"

// ResourceWidgetDependency is holding the info about a javascript script to be included.
func ResourceWidgetDependency() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_widget_dependency` manages JS and CSS includes for a Widget configuration within ServiceNow.",

		CreateContext: createResourceWidgetDependency,
		ReadContext:   readResourceWidgetDependency,
		UpdateContext: updateResourceWidgetDependency,
		DeleteContext: deleteResourceWidgetDependency,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			widgetDependencyName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the widget dependency.",
			},
			widgetDependencyModule: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The AngularJS module name for this dependency.",
			},
			widgetDependencyPageLoad: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this dependency is loaded on every page load.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceWidgetDependency(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	widgetDependency := &client.WidgetDependency{}
	if err := snowClient.GetObject(client.EndpointWidgetDependency, data.Id(), widgetDependency); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromWidgetDependency(data, widgetDependency)

	return nil
}

func createResourceWidgetDependency(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	widgetDependency := resourceToWidgetDependency(data)
	if err := snowClient.CreateObject(client.EndpointWidgetDependency, widgetDependency); err != nil {
		return diag.FromErr(err)
	}

	resourceFromWidgetDependency(data, widgetDependency)

	return readResourceWidgetDependency(ctx, data, serviceNowClient)
}

func updateResourceWidgetDependency(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointWidgetDependency, resourceToWidgetDependency(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceWidgetDependency(ctx, data, serviceNowClient)
}

func deleteResourceWidgetDependency(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointWidgetDependency, data.Id()))
}

func resourceFromWidgetDependency(data *schema.ResourceData, widgetDependency *client.WidgetDependency) {
	data.SetId(widgetDependency.ID)
	data.Set(widgetDependencyName, widgetDependency.Name)
	data.Set(widgetDependencyModule, widgetDependency.Module)
	data.Set(widgetDependencyPageLoad, widgetDependency.PageLoad)
	data.Set(commonScope, widgetDependency.Scope)
}

func resourceToWidgetDependency(data *schema.ResourceData) *client.WidgetDependency {
	widgetDependency := client.WidgetDependency{
		Name:     data.Get(widgetDependencyName).(string),
		Module:   data.Get(widgetDependencyModule).(string),
		PageLoad: data.Get(widgetDependencyPageLoad).(bool),
	}
	widgetDependency.ID = data.Id()
	widgetDependency.Scope = data.Get(commonScope).(string)
	return &widgetDependency
}
