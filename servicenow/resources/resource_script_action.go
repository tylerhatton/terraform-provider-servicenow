package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const scriptActionName = "name"
const scriptActionEventName = "event_name"
const scriptActionScript = "script"
const scriptActionActive = "active"
const scriptActionDescription = "description"
const scriptActionOrder = "order"
const scriptActionCondition = "condition"
const scriptActionSynchronous = "synchronous"

// ResourceScriptAction manages an event-driven script action in ServiceNow.
func ResourceScriptAction() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_script_action` manages an event-driven script action (sysevent_script_action) within ServiceNow.",

		CreateContext: createResourceScriptAction,
		ReadContext:   readResourceScriptAction,
		UpdateContext: updateResourceScriptAction,
		DeleteContext: deleteResourceScriptAction,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			scriptActionName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the script action.",
			},
			scriptActionEventName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the event that this script action fires on.",
			},
			scriptActionScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Javascript code to execute when the event fires.",
			},
			scriptActionActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this script action is enabled.",
			},
			scriptActionDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the script action and its purpose.",
			},
			scriptActionOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Execution order for this script action.",
			},
			scriptActionCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Condition script that must evaluate true for the action to fire.",
			},
			scriptActionSynchronous: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, run the script action synchronously.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceScriptAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scriptAction := &client.ScriptAction{}
	if err := snowClient.GetObject(ctx, client.EndpointScriptAction, data.Id(), scriptAction); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromScriptAction(data, scriptAction)

	return nil
}

func createResourceScriptAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scriptAction := resourceToScriptAction(data)
	if err := snowClient.CreateObject(ctx, client.EndpointScriptAction, scriptAction); err != nil {
		return diag.FromErr(err)
	}

	resourceFromScriptAction(data, scriptAction)

	return readResourceScriptAction(ctx, data, serviceNowClient)
}

func updateResourceScriptAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointScriptAction, resourceToScriptAction(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceScriptAction(ctx, data, serviceNowClient)
}

func deleteResourceScriptAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointScriptAction, data.Id()))
}

func resourceFromScriptAction(data *schema.ResourceData, scriptAction *client.ScriptAction) {
	data.SetId(scriptAction.ID)
	data.Set(scriptActionName, scriptAction.Name)
	data.Set(scriptActionEventName, scriptAction.EventName)
	data.Set(scriptActionScript, scriptAction.Script)
	data.Set(scriptActionActive, scriptAction.Active)
	data.Set(scriptActionDescription, scriptAction.Description)
	data.Set(scriptActionOrder, scriptAction.Order)
	data.Set(scriptActionCondition, scriptAction.Condition)
	data.Set(scriptActionSynchronous, scriptAction.Synchronous)
	data.Set(commonProtectionPolicy, scriptAction.ProtectionPolicy)
	data.Set(commonScope, scriptAction.Scope)
}

func resourceToScriptAction(data *schema.ResourceData) *client.ScriptAction {
	scriptAction := client.ScriptAction{
		Name:        data.Get(scriptActionName).(string),
		EventName:   data.Get(scriptActionEventName).(string),
		Script:      data.Get(scriptActionScript).(string),
		Active:      data.Get(scriptActionActive).(bool),
		Description: data.Get(scriptActionDescription).(string),
		Order:       data.Get(scriptActionOrder).(int),
		Condition:   data.Get(scriptActionCondition).(string),
		Synchronous: data.Get(scriptActionSynchronous).(bool),
	}
	scriptAction.ID = data.Id()
	scriptAction.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	scriptAction.Scope = data.Get(commonScope).(string)
	return &scriptAction
}
