package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const jsIncludeSource = "source"
const jsIncludeDisplayName = "display_name"
const jsIncludeURL = "url"
const jsIncludeUIScriptID = "ui_script_id"

// ResourceJsInclude is holding the info about a javascript script to be included.
func ResourceJsInclude() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_js_include` manages a javascript script within ServiceNow.",

		CreateContext: createResourceJsInclude,
		ReadContext:   readResourceJsInclude,
		UpdateContext: updateResourceJsInclude,
		DeleteContext: deleteResourceJsInclude,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			jsIncludeSource: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "url",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					warns, errs = validateStringValue(val.(string), key, []string{"url", "local"})
					return
				},
				Description: "Source type of the JS include. Can be 'url' for an external link or 'local' for a UI script.",
			},
			jsIncludeDisplayName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the JS include.",
			},
			jsIncludeURL: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "URL of the external JavaScript file when source is set to 'url'.",
			},
			jsIncludeUIScriptID: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The Sys ID of the UI Script to include when source is set to 'local'.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceJsInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	jsInclude := &client.JsInclude{}
	if err := snowClient.GetObject(client.EndpointJsInclude, data.Id(), jsInclude); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromJsInclude(data, jsInclude)

	return nil
}

func createResourceJsInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	jsInclude := resourceToJsInclude(data)
	if err := snowClient.CreateObject(client.EndpointJsInclude, jsInclude); err != nil {
		return diag.FromErr(err)
	}

	resourceFromJsInclude(data, jsInclude)

	return readResourceJsInclude(ctx, data, serviceNowClient)
}

func updateResourceJsInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointJsInclude, resourceToJsInclude(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceJsInclude(ctx, data, serviceNowClient)
}

func deleteResourceJsInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointJsInclude, data.Id()))
}

func resourceFromJsInclude(data *schema.ResourceData, jsInclude *client.JsInclude) {
	data.SetId(jsInclude.ID)
	data.Set(jsIncludeSource, jsInclude.Source)
	data.Set(jsIncludeDisplayName, jsInclude.DisplayName)
	data.Set(jsIncludeURL, jsInclude.URL)
	data.Set(jsIncludeUIScriptID, jsInclude.UIScriptID)
	data.Set(commonScope, jsInclude.Scope)
}

func resourceToJsInclude(data *schema.ResourceData) *client.JsInclude {
	jsInclude := client.JsInclude{
		Source:      data.Get(jsIncludeSource).(string),
		DisplayName: data.Get(jsIncludeDisplayName).(string),
		URL:         data.Get(jsIncludeURL).(string),
		UIScriptID:  data.Get(jsIncludeUIScriptID).(string),
	}
	jsInclude.ID = data.Id()
	jsInclude.Scope = data.Get(commonScope).(string)
	return &jsInclude
}
