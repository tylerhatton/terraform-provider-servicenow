package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const uiScriptName = "name"
const uiScriptDescription = "description"
const uiScriptScript = "script"
const uiScriptActive = "active"
const uiScriptUIType = "type"
const uiScriptAPIName = "api_name"

// ResourceUIScript manages a UI Script in ServiceNow which can be added to any other UI component.
func ResourceUIScript() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_ui_script` manages a UI Script configuration within ServiceNow.",

		CreateContext: createResourceUIScript,
		ReadContext:   readResourceUIScript,
		UpdateContext: updateResourceUIScript,
		DeleteContext: deleteResourceUIScript,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			uiScriptName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the UI Script.",
			},
			uiScriptDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the UI Script.",
			},
			uiScriptScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The JavaScript body of the UI Script.",
			},
			uiScriptActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', this UI Script is enabled and available for use.",
			},
			uiScriptUIType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "all",
				Description: "The UI type this script applies to. Valid values are 'all', 'desktop', or 'mobile'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					warns, errs = validateStringValue(val.(string), key, []string{"all", "desktop", "mobile"})
					return
				},
			},
			uiScriptAPIName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scoped API name of the UI Script.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceUIScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiScript := &client.UIScript{}
	if err := snowClient.GetObject(client.EndpointUIScript, data.Id(), uiScript); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIScript(data, uiScript)

	return nil
}

func createResourceUIScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiScript := resourceToUIScript(data)
	if err := snowClient.CreateObject(client.EndpointUIScript, uiScript); err != nil {
		return diag.FromErr(err)
	}

	resourceFromUIScript(data, uiScript)

	return readResourceUIScript(ctx, data, serviceNowClient)
}

func updateResourceUIScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointUIScript, resourceToUIScript(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceUIScript(ctx, data, serviceNowClient)
}

func deleteResourceUIScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointUIScript, data.Id()))
}

func resourceFromUIScript(data *schema.ResourceData, script *client.UIScript) {
	var typeString string
	switch script.UIType {
	case "1":
		typeString = "mobile"
	case "0":
		typeString = "desktop"
	default:
		typeString = "all"
	}

	data.SetId(script.ID)
	data.Set(uiScriptName, script.Name)
	data.Set(uiScriptDescription, script.Description)
	data.Set(uiScriptScript, script.Script)
	data.Set(uiScriptActive, script.Active)
	data.Set(uiScriptUIType, typeString)
	data.Set(uiScriptAPIName, script.APIName)
}

func resourceToUIScript(data *schema.ResourceData) *client.UIScript {
	var typeInt string
	switch data.Get(uiScriptUIType).(string) {
	case "mobile":
		typeInt = "1"
	case "desktop":
		typeInt = "0"
	default:
		typeInt = "10"
	}

	uiScript := client.UIScript{
		Name:        data.Get(uiScriptName).(string),
		Description: data.Get(uiScriptDescription).(string),
		Script:      data.Get(uiScriptScript).(string),
		Active:      data.Get(uiScriptActive).(bool),
		UIType:      typeInt,
	}
	uiScript.ID = data.Id()
	uiScript.Scope = data.Get(commonScope).(string)
	return &uiScript
}
