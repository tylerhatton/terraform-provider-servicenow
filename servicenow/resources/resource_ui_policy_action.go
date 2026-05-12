package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const uiPolicyActionUIPolicy = "ui_policy"
const uiPolicyActionFieldName = "field_name"
const uiPolicyActionVisible = "visible"
const uiPolicyActionMandatory = "mandatory"
const uiPolicyActionDisabled = "disabled"
const uiPolicyActionCleared = "cleared"

// ResourceUIPolicyAction manages a UI Policy Action in ServiceNow.
func ResourceUIPolicyAction() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_ui_policy_action` manages a UI policy action within ServiceNow. UI policy actions are child records of a UI policy that specify which fields are affected and how (visible, mandatory, disabled).",

		CreateContext: createResourceUIPolicyAction,
		ReadContext:   readResourceUIPolicyAction,
		UpdateContext: updateResourceUIPolicyAction,
		DeleteContext: deleteResourceUIPolicyAction,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			uiPolicyActionUIPolicy: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sys ID of the parent UI policy this action belongs to. Cannot be changed once created.",
			},
			uiPolicyActionFieldName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the field this action targets.",
			},
			uiPolicyActionVisible: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ignore",
				Description:  "Whether to set the field visible. Valid values are 'true', 'false', or 'ignore'.",
				ValidateFunc: validation.StringInSlice([]string{"true", "false", "ignore"}, false),
			},
			uiPolicyActionMandatory: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ignore",
				Description:  "Whether to set the field mandatory. Valid values are 'true', 'false', or 'ignore'.",
				ValidateFunc: validation.StringInSlice([]string{"true", "false", "ignore"}, false),
			},
			uiPolicyActionDisabled: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ignore",
				Description:  "Whether to set the field disabled (read-only). Valid values are 'true', 'false', or 'ignore'.",
				ValidateFunc: validation.StringInSlice([]string{"true", "false", "ignore"}, false),
			},
			uiPolicyActionCleared: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the field's value is cleared when the policy condition is met.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceUIPolicyAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiPolicyAction := &client.UIPolicyAction{}
	if err := snowClient.GetObject(ctx, client.EndpointUIPolicyAction, data.Id(), uiPolicyAction); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIPolicyAction(data, uiPolicyAction)

	return nil
}

func createResourceUIPolicyAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiPolicyAction := resourceToUIPolicyAction(data)
	if err := snowClient.CreateObject(ctx, client.EndpointUIPolicyAction, uiPolicyAction); err != nil {
		return diag.FromErr(err)
	}

	resourceFromUIPolicyAction(data, uiPolicyAction)

	return readResourceUIPolicyAction(ctx, data, serviceNowClient)
}

func updateResourceUIPolicyAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointUIPolicyAction, resourceToUIPolicyAction(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceUIPolicyAction(ctx, data, serviceNowClient)
}

func deleteResourceUIPolicyAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointUIPolicyAction, data.Id()))
}

func resourceFromUIPolicyAction(data *schema.ResourceData, uiPolicyAction *client.UIPolicyAction) {
	data.SetId(uiPolicyAction.ID)
	data.Set(uiPolicyActionUIPolicy, uiPolicyAction.UIPolicy)
	data.Set(uiPolicyActionFieldName, uiPolicyAction.FieldName)
	data.Set(uiPolicyActionVisible, uiPolicyAction.Visible)
	data.Set(uiPolicyActionMandatory, uiPolicyAction.Mandatory)
	data.Set(uiPolicyActionDisabled, uiPolicyAction.Disabled)
	data.Set(uiPolicyActionCleared, uiPolicyAction.Cleared)
	data.Set(commonScope, uiPolicyAction.Scope)
}

func resourceToUIPolicyAction(data *schema.ResourceData) *client.UIPolicyAction {
	uiPolicyAction := client.UIPolicyAction{
		UIPolicy:  data.Get(uiPolicyActionUIPolicy).(string),
		FieldName: data.Get(uiPolicyActionFieldName).(string),
		Visible:   data.Get(uiPolicyActionVisible).(string),
		Mandatory: data.Get(uiPolicyActionMandatory).(string),
		Disabled:  data.Get(uiPolicyActionDisabled).(string),
		Cleared:   data.Get(uiPolicyActionCleared).(bool),
	}
	uiPolicyAction.ID = data.Id()
	uiPolicyAction.Scope = data.Get(commonScope).(string)
	return &uiPolicyAction
}
