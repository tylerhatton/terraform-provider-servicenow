package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const extensionPointName = "name"
const extensionPointDescription = "description"
const extensionPointRestrictScope = "restrict_scope"
const extensionPointExample = "example"
const extensionPointAPIName = "api_name"

// ResourceExtensionPoint is holding the info about a scripted extension point.
func ResourceExtensionPoint() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_extension_point` manages a scripted extension point within ServiceNow.",

		CreateContext: createResourceExtensionPoint,
		ReadContext:   readResourceExtensionPoint,
		UpdateContext: updateResourceExtensionPoint,
		DeleteContext: deleteResourceExtensionPoint,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			extensionPointName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name of the extension point.",
			},
			extensionPointDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the required implementation of the extension point.",
			},
			extensionPointExample: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Example implementation code.",
			},
			extensionPointRestrictScope: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Only allow extension instances within this point's scope.",
			},
			extensionPointAPIName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the extension point API, that is pre-pended with the application scope to which it applies. This is a system-assigned name and cannot be changed.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceExtensionPoint(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	extensionPoint := &client.ExtensionPoint{}
	if err := snowClient.GetObject(client.EndpointExtensionPoint, data.Id(), extensionPoint); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromExtensionPoint(data, extensionPoint)

	return nil
}

func createResourceExtensionPoint(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	extensionPoint := resourceToExtensionPoint(data)
	if err := snowClient.CreateObject(client.EndpointExtensionPoint, extensionPoint); err != nil {
		return diag.FromErr(err)
	}

	resourceFromExtensionPoint(data, extensionPoint)

	return readResourceExtensionPoint(ctx, data, serviceNowClient)
}

func updateResourceExtensionPoint(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointExtensionPoint, resourceToExtensionPoint(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceExtensionPoint(ctx, data, serviceNowClient)
}

func deleteResourceExtensionPoint(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointExtensionPoint, data.Id()))
}

func resourceFromExtensionPoint(data *schema.ResourceData, extensionPoint *client.ExtensionPoint) {
	data.SetId(extensionPoint.ID)
	data.Set(extensionPointName, extensionPoint.Name)
	data.Set(extensionPointDescription, extensionPoint.Description)
	data.Set(extensionPointRestrictScope, extensionPoint.RestrictScope)
	data.Set(extensionPointExample, extensionPoint.Example)
	data.Set(extensionPointAPIName, extensionPoint.APIName)
	data.Set(commonScope, extensionPoint.Scope)
}

func resourceToExtensionPoint(data *schema.ResourceData) *client.ExtensionPoint {
	extensionPoint := client.ExtensionPoint{
		Name:          data.Get(extensionPointName).(string),
		Description:   data.Get(extensionPointDescription).(string),
		RestrictScope: data.Get(extensionPointRestrictScope).(bool),
		Example:       data.Get(extensionPointExample).(string),
	}
	extensionPoint.ID = data.Id()
	extensionPoint.Scope = data.Get(commonScope).(string)
	return &extensionPoint
}
