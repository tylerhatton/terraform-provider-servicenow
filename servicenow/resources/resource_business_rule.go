package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const businessRuleName = "name"
const businessRuleTable = "table"
const businessRuleWhen = "when"
const businessRuleOrder = "order"
const businessRuleActive = "active"
const businessRuleCondition = "condition"
const businessRuleFilterCondition = "filter_condition"
const businessRuleScript = "script"
const businessRuleDescription = "description"
const businessRuleAdvanced = "advanced"
const businessRuleActionInsert = "action_insert"
const businessRuleActionUpdate = "action_update"
const businessRuleActionDelete = "action_delete"
const businessRuleActionQuery = "action_query"
const businessRuleRoleConditions = "role_conditions"
const businessRulePriority = "priority"
const businessRuleAbortAction = "abort_action"
const businessRuleAddMessage = "add_message"
const businessRuleIsRest = "is_rest"
const businessRuleClientCallable = "client_callable"

// ResourceBusinessRule manages a Business Rule in ServiceNow.
func ResourceBusinessRule() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_business_rule` manages a business rule within ServiceNow. Business rules execute server-side JavaScript when records are inserted, updated, deleted, queried, or displayed.",

		CreateContext: createResourceBusinessRule,
		ReadContext:   readResourceBusinessRule,
		UpdateContext: updateResourceBusinessRule,
		DeleteContext: deleteResourceBusinessRule,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			businessRuleName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the business rule.",
			},
			businessRuleTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the table this business rule fires on. Cannot be changed once created.",
			},
			businessRuleWhen: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "before",
				Description:  "When the business rule should execute. Valid values are 'before', 'after', 'async', or 'display'.",
				ValidateFunc: validation.StringInSlice([]string{"before", "after", "async", "display"}, false),
			},
			businessRuleOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The order in which this business rule executes relative to others on the same table.",
			},
			businessRuleActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this business rule is enabled.",
			},
			businessRuleCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The condition under which this business rule executes.",
			},
			businessRuleFilterCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "An encoded query filter to determine when this business rule applies.",
			},
			businessRuleScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The JavaScript to execute when the business rule fires.",
			},
			businessRuleDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of what this business rule does.",
			},
			businessRuleAdvanced: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If true, the rule uses advanced (script-based) configuration. ServiceNow auto-sets this to true when a script is supplied.",
			},
			businessRuleActionInsert: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the business rule fires on insert operations.",
			},
			businessRuleActionUpdate: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the business rule fires on update operations.",
			},
			businessRuleActionDelete: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the business rule fires on delete operations.",
			},
			businessRuleActionQuery: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the business rule fires on query operations.",
			},
			businessRuleRoleConditions: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-separated list of roles required to trigger this business rule.",
			},
			businessRulePriority: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The priority of this business rule.",
			},
			businessRuleAbortAction: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, aborts the current action when the rule's condition is met.",
			},
			businessRuleAddMessage: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, an informational message is shown when the rule fires.",
			},
			businessRuleIsRest: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this business rule applies to REST API operations.",
			},
			businessRuleClientCallable: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this business rule can be called from client-side scripts.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceBusinessRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	businessRule := &client.BusinessRule{}
	if err := snowClient.GetObject(ctx, client.EndpointBusinessRule, data.Id(), businessRule); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromBusinessRule(data, businessRule)

	return nil
}

func createResourceBusinessRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	businessRule := resourceToBusinessRule(data)
	if err := snowClient.CreateObject(ctx, client.EndpointBusinessRule, businessRule); err != nil {
		return diag.FromErr(err)
	}

	resourceFromBusinessRule(data, businessRule)

	return readResourceBusinessRule(ctx, data, serviceNowClient)
}

func updateResourceBusinessRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointBusinessRule, resourceToBusinessRule(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceBusinessRule(ctx, data, serviceNowClient)
}

func deleteResourceBusinessRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointBusinessRule, data.Id()))
}

func resourceFromBusinessRule(data *schema.ResourceData, businessRule *client.BusinessRule) {
	data.SetId(businessRule.ID)
	data.Set(businessRuleName, businessRule.Name)
	data.Set(businessRuleTable, businessRule.Table)
	data.Set(businessRuleWhen, businessRule.When)
	data.Set(businessRuleOrder, businessRule.Order)
	data.Set(businessRuleActive, businessRule.Active)
	data.Set(businessRuleCondition, businessRule.Condition)
	data.Set(businessRuleFilterCondition, businessRule.FilterCondition)
	data.Set(businessRuleScript, businessRule.Script)
	data.Set(businessRuleDescription, businessRule.Description)
	data.Set(businessRuleAdvanced, businessRule.Advanced)
	data.Set(businessRuleActionInsert, businessRule.ActionInsert)
	data.Set(businessRuleActionUpdate, businessRule.ActionUpdate)
	data.Set(businessRuleActionDelete, businessRule.ActionDelete)
	data.Set(businessRuleActionQuery, businessRule.ActionQuery)
	data.Set(businessRuleRoleConditions, businessRule.RoleConditions)
	data.Set(businessRulePriority, businessRule.Priority)
	data.Set(businessRuleAbortAction, businessRule.AbortAction)
	data.Set(businessRuleAddMessage, businessRule.AddMessage)
	data.Set(businessRuleIsRest, businessRule.IsRest)
	data.Set(businessRuleClientCallable, businessRule.ClientCallable)
	data.Set(commonProtectionPolicy, businessRule.ProtectionPolicy)
	data.Set(commonScope, businessRule.Scope)
}

func resourceToBusinessRule(data *schema.ResourceData) *client.BusinessRule {
	businessRule := client.BusinessRule{
		Name:            data.Get(businessRuleName).(string),
		Table:           data.Get(businessRuleTable).(string),
		When:            data.Get(businessRuleWhen).(string),
		Order:           data.Get(businessRuleOrder).(int),
		Active:          data.Get(businessRuleActive).(bool),
		Condition:       data.Get(businessRuleCondition).(string),
		FilterCondition: data.Get(businessRuleFilterCondition).(string),
		Script:          data.Get(businessRuleScript).(string),
		Description:     data.Get(businessRuleDescription).(string),
		Advanced:        data.Get(businessRuleAdvanced).(bool),
		ActionInsert:    data.Get(businessRuleActionInsert).(bool),
		ActionUpdate:    data.Get(businessRuleActionUpdate).(bool),
		ActionDelete:    data.Get(businessRuleActionDelete).(bool),
		ActionQuery:     data.Get(businessRuleActionQuery).(bool),
		RoleConditions:  data.Get(businessRuleRoleConditions).(string),
		Priority:        data.Get(businessRulePriority).(int),
		AbortAction:     data.Get(businessRuleAbortAction).(bool),
		AddMessage:      data.Get(businessRuleAddMessage).(bool),
		IsRest:          data.Get(businessRuleIsRest).(bool),
		ClientCallable:  data.Get(businessRuleClientCallable).(bool),
	}
	businessRule.ID = data.Id()
	businessRule.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	businessRule.Scope = data.Get(commonScope).(string)
	return &businessRule
}
