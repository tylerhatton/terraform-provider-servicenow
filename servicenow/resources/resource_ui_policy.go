package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const uiPolicyShortDescription = "short_description"
const uiPolicyTable = "table"
const uiPolicyActive = "active"
const uiPolicyConditions = "conditions"
const uiPolicyReverseIfFalse = "reverse_if_false"
const uiPolicyRunScripts = "run_scripts"
const uiPolicyScriptTrue = "script_true"
const uiPolicyScriptFalse = "script_false"
const uiPolicyOnLoad = "on_load"
const uiPolicyGlobal = "global"
const uiPolicyOrder = "order"
const uiPolicyView = "view"
const uiPolicyDescription = "description"

// ResourceUIPolicy manages a UI Policy in ServiceNow.
func ResourceUIPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_ui_policy` manages a UI policy within ServiceNow. UI policies dynamically change form fields (visibility, mandatory, read-only) based on conditions.",

		CreateContext: createResourceUIPolicy,
		ReadContext:   readResourceUIPolicy,
		UpdateContext: updateResourceUIPolicy,
		DeleteContext: deleteResourceUIPolicy,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			uiPolicyShortDescription: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Short description identifying this UI policy.",
			},
			uiPolicyTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the table this UI policy applies to. Cannot be changed once created.",
			},
			uiPolicyActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this UI policy is enabled.",
			},
			uiPolicyConditions: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Encoded query conditions that trigger this UI policy.",
			},
			uiPolicyReverseIfFalse: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the actions are reversed when the condition is false.",
			},
			uiPolicyRunScripts: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, run the script_true and script_false scripts when the policy evaluates.",
			},
			uiPolicyScriptTrue: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script to execute when the condition evaluates to true. ServiceNow seeds a default function body when the policy is created; specify a value here to override it.",
			},
			uiPolicyScriptFalse: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script to execute when the condition evaluates to false. ServiceNow seeds a default function body when the policy is created; specify a value here to override it.",
			},
			uiPolicyOnLoad: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the UI policy runs when the form loads.",
			},
			uiPolicyGlobal: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the UI policy applies to all views of the table.",
			},
			uiPolicyOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The order in which this UI policy evaluates relative to others.",
			},
			uiPolicyView: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the sys_ui_view this UI policy applies to.",
			},
			uiPolicyDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of what this UI policy does.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceUIPolicy(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiPolicy := &client.UIPolicy{}
	if err := snowClient.GetObject(ctx, client.EndpointUIPolicy, data.Id(), uiPolicy); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIPolicy(data, uiPolicy)

	return nil
}

func createResourceUIPolicy(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiPolicy := resourceToUIPolicy(data)
	if err := snowClient.CreateObject(ctx, client.EndpointUIPolicy, uiPolicy); err != nil {
		return diag.FromErr(err)
	}

	resourceFromUIPolicy(data, uiPolicy)

	return readResourceUIPolicy(ctx, data, serviceNowClient)
}

func updateResourceUIPolicy(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointUIPolicy, resourceToUIPolicy(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceUIPolicy(ctx, data, serviceNowClient)
}

func deleteResourceUIPolicy(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointUIPolicy, data.Id()))
}

func resourceFromUIPolicy(data *schema.ResourceData, uiPolicy *client.UIPolicy) {
	data.SetId(uiPolicy.ID)
	data.Set(uiPolicyShortDescription, uiPolicy.ShortDescription)
	data.Set(uiPolicyTable, uiPolicy.Table)
	data.Set(uiPolicyActive, uiPolicy.Active)
	data.Set(uiPolicyConditions, uiPolicy.Conditions)
	data.Set(uiPolicyReverseIfFalse, uiPolicy.ReverseIfFalse)
	data.Set(uiPolicyRunScripts, uiPolicy.RunScripts)
	data.Set(uiPolicyScriptTrue, uiPolicy.ScriptTrue)
	data.Set(uiPolicyScriptFalse, uiPolicy.ScriptFalse)
	data.Set(uiPolicyOnLoad, uiPolicy.OnLoad)
	data.Set(uiPolicyGlobal, uiPolicy.Global)
	data.Set(uiPolicyOrder, uiPolicy.Order)
	data.Set(uiPolicyView, uiPolicy.View)
	data.Set(uiPolicyDescription, uiPolicy.Description)
	data.Set(commonProtectionPolicy, uiPolicy.ProtectionPolicy)
	data.Set(commonScope, uiPolicy.Scope)
}

func resourceToUIPolicy(data *schema.ResourceData) *client.UIPolicy {
	uiPolicy := client.UIPolicy{
		ShortDescription: data.Get(uiPolicyShortDescription).(string),
		Table:            data.Get(uiPolicyTable).(string),
		Active:           data.Get(uiPolicyActive).(bool),
		Conditions:       data.Get(uiPolicyConditions).(string),
		ReverseIfFalse:   data.Get(uiPolicyReverseIfFalse).(bool),
		RunScripts:       data.Get(uiPolicyRunScripts).(bool),
		ScriptTrue:       data.Get(uiPolicyScriptTrue).(string),
		ScriptFalse:      data.Get(uiPolicyScriptFalse).(string),
		OnLoad:           data.Get(uiPolicyOnLoad).(bool),
		Global:           data.Get(uiPolicyGlobal).(bool),
		Order:            data.Get(uiPolicyOrder).(int),
		View:             data.Get(uiPolicyView).(string),
		Description:      data.Get(uiPolicyDescription).(string),
	}
	uiPolicy.ID = data.Id()
	uiPolicy.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	uiPolicy.Scope = data.Get(commonScope).(string)
	return &uiPolicy
}
