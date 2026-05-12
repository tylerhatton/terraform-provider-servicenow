package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const clientScriptName = "name"
const clientScriptTable = "table"
const clientScriptType = "type"
const clientScriptFieldName = "field_name"
const clientScriptScript = "script"
const clientScriptActive = "active"
const clientScriptDescription = "description"
const clientScriptGlobal = "global"
const clientScriptAppliesExtended = "applies_extended"
const clientScriptCondition = "condition"
const clientScriptMessages = "messages"
const clientScriptView = "view"
const clientScriptUIType = "ui_type"

// ResourceClientScript manages a Client Script in ServiceNow.
func ResourceClientScript() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_client_script` manages a client script within ServiceNow. Client scripts run JavaScript in the browser when forms are loaded, changed, submitted, or when cells are edited in lists.",

		CreateContext: createResourceClientScript,
		ReadContext:   readResourceClientScript,
		UpdateContext: updateResourceClientScript,
		DeleteContext: deleteResourceClientScript,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			clientScriptName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the client script.",
			},
			clientScriptTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the table this client script applies to. Cannot be changed once created.",
			},
			clientScriptType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "onLoad",
				Description:  "When the client script runs. Valid values are 'onLoad', 'onChange', 'onSubmit', or 'onCellEdit'.",
				ValidateFunc: validation.StringInSlice([]string{"onLoad", "onChange", "onSubmit", "onCellEdit"}, false),
			},
			clientScriptFieldName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Field name this client script applies to. Required when type is 'onChange' or 'onCellEdit'.",
			},
			clientScriptScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The JavaScript to execute on the client.",
			},
			clientScriptActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this client script is enabled.",
			},
			clientScriptDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of what this client script does.",
			},
			clientScriptGlobal: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the script applies to subtables of the specified table.",
			},
			clientScriptAppliesExtended: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the script applies to extended tables.",
			},
			clientScriptCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Condition under which this client script runs.",
			},
			clientScriptMessages: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Messages that can be referenced by this client script.",
			},
			clientScriptView: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the sys_ui_view this client script applies to.",
			},
			clientScriptUIType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "0",
				Description:  "The UI type this script applies to. '0' for desktop, '1' for mobile, '10' for both.",
				ValidateFunc: validation.StringInSlice([]string{"0", "1", "10"}, false),
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceClientScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	clientScript := &client.ClientScript{}
	if err := snowClient.GetObject(ctx, client.EndpointClientScript, data.Id(), clientScript); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromClientScript(data, clientScript)

	return nil
}

func createResourceClientScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	clientScript := resourceToClientScript(data)
	if err := snowClient.CreateObject(ctx, client.EndpointClientScript, clientScript); err != nil {
		return diag.FromErr(err)
	}

	resourceFromClientScript(data, clientScript)

	return readResourceClientScript(ctx, data, serviceNowClient)
}

func updateResourceClientScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointClientScript, resourceToClientScript(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceClientScript(ctx, data, serviceNowClient)
}

func deleteResourceClientScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointClientScript, data.Id()))
}

func resourceFromClientScript(data *schema.ResourceData, clientScript *client.ClientScript) {
	data.SetId(clientScript.ID)
	data.Set(clientScriptName, clientScript.Name)
	data.Set(clientScriptTable, clientScript.Table)
	data.Set(clientScriptType, clientScript.Type)
	data.Set(clientScriptFieldName, clientScript.FieldName)
	data.Set(clientScriptScript, clientScript.Script)
	data.Set(clientScriptActive, clientScript.Active)
	data.Set(clientScriptDescription, clientScript.Description)
	data.Set(clientScriptGlobal, clientScript.Global)
	data.Set(clientScriptAppliesExtended, clientScript.AppliesExtended)
	data.Set(clientScriptCondition, clientScript.Condition)
	data.Set(clientScriptMessages, clientScript.Messages)
	data.Set(clientScriptView, clientScript.View)
	data.Set(clientScriptUIType, clientScript.UIType)
	data.Set(commonProtectionPolicy, clientScript.ProtectionPolicy)
	data.Set(commonScope, clientScript.Scope)
}

func resourceToClientScript(data *schema.ResourceData) *client.ClientScript {
	clientScript := client.ClientScript{
		Name:            data.Get(clientScriptName).(string),
		Table:           data.Get(clientScriptTable).(string),
		Type:            data.Get(clientScriptType).(string),
		FieldName:       data.Get(clientScriptFieldName).(string),
		Script:          data.Get(clientScriptScript).(string),
		Active:          data.Get(clientScriptActive).(bool),
		Description:     data.Get(clientScriptDescription).(string),
		Global:          data.Get(clientScriptGlobal).(bool),
		AppliesExtended: data.Get(clientScriptAppliesExtended).(bool),
		Condition:       data.Get(clientScriptCondition).(string),
		Messages:        data.Get(clientScriptMessages).(string),
		View:            data.Get(clientScriptView).(string),
		UIType:          data.Get(clientScriptUIType).(string),
	}
	clientScript.ID = data.Id()
	clientScript.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	clientScript.Scope = data.Get(commonScope).(string)
	return &clientScript
}
