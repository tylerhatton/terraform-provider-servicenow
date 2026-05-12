package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const assignmentRuleName = "name"
const assignmentRuleTable = "table"
const assignmentRuleActive = "active"
const assignmentRuleOrder = "order"
const assignmentRuleDescription = "description"
const assignmentRuleCondition = "condition"
const assignmentRuleGroup = "assignment_group"
const assignmentRuleUser = "user"
const assignmentRuleScript = "script"
const assignmentRuleMatchFor = "match_for"

// ResourceAssignmentRule manages an assignment rule in ServiceNow.
func ResourceAssignmentRule() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_assignment_rule` manages an assignment rule (sysrule_assignment) within ServiceNow.",

		CreateContext: createResourceAssignmentRule,
		ReadContext:   readResourceAssignmentRule,
		UpdateContext: updateResourceAssignmentRule,
		DeleteContext: deleteResourceAssignmentRule,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			assignmentRuleName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the assignment rule.",
			},
			assignmentRuleTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Table that this assignment rule applies to.",
			},
			assignmentRuleActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this assignment rule is enabled.",
			},
			assignmentRuleOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Order in which the rule is evaluated.",
			},
			assignmentRuleDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the assignment rule.",
			},
			assignmentRuleCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Encoded query condition that must match for the rule to apply.",
			},
			assignmentRuleGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "sys_id of the group to assign the record to.",
			},
			assignmentRuleUser: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "sys_id of the user to assign the record to.",
			},
			assignmentRuleScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script to determine the assignment. ServiceNow seeds an example script body when this is omitted.",
			},
			assignmentRuleMatchFor: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Match conditions setting (e.g. ALL). ServiceNow defaults this to ALL when omitted.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceAssignmentRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	assignmentRule := &client.AssignmentRule{}
	if err := snowClient.GetObject(ctx, client.EndpointAssignmentRule, data.Id(), assignmentRule); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromAssignmentRule(data, assignmentRule)

	return nil
}

func createResourceAssignmentRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	assignmentRule := resourceToAssignmentRule(data)
	if err := snowClient.CreateObject(ctx, client.EndpointAssignmentRule, assignmentRule); err != nil {
		return diag.FromErr(err)
	}

	resourceFromAssignmentRule(data, assignmentRule)

	return readResourceAssignmentRule(ctx, data, serviceNowClient)
}

func updateResourceAssignmentRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointAssignmentRule, resourceToAssignmentRule(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceAssignmentRule(ctx, data, serviceNowClient)
}

func deleteResourceAssignmentRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointAssignmentRule, data.Id()))
}

func resourceFromAssignmentRule(data *schema.ResourceData, assignmentRule *client.AssignmentRule) {
	data.SetId(assignmentRule.ID)
	data.Set(assignmentRuleName, assignmentRule.Name)
	data.Set(assignmentRuleTable, assignmentRule.Table)
	data.Set(assignmentRuleActive, assignmentRule.Active)
	data.Set(assignmentRuleOrder, assignmentRule.Order)
	data.Set(assignmentRuleDescription, assignmentRule.Description)
	data.Set(assignmentRuleCondition, assignmentRule.Condition)
	data.Set(assignmentRuleGroup, assignmentRule.AssignmentGroup)
	data.Set(assignmentRuleUser, assignmentRule.User)
	data.Set(assignmentRuleScript, assignmentRule.Script)
	data.Set(assignmentRuleMatchFor, assignmentRule.MatchFor)
	data.Set(commonProtectionPolicy, assignmentRule.ProtectionPolicy)
	data.Set(commonScope, assignmentRule.Scope)
}

func resourceToAssignmentRule(data *schema.ResourceData) *client.AssignmentRule {
	assignmentRule := client.AssignmentRule{
		Name:            data.Get(assignmentRuleName).(string),
		Table:           data.Get(assignmentRuleTable).(string),
		Active:          data.Get(assignmentRuleActive).(bool),
		Order:           data.Get(assignmentRuleOrder).(int),
		Description:     data.Get(assignmentRuleDescription).(string),
		Condition:       data.Get(assignmentRuleCondition).(string),
		AssignmentGroup: data.Get(assignmentRuleGroup).(string),
		User:            data.Get(assignmentRuleUser).(string),
		Script:          data.Get(assignmentRuleScript).(string),
		MatchFor:        data.Get(assignmentRuleMatchFor).(string),
	}
	assignmentRule.ID = data.Id()
	assignmentRule.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	assignmentRule.Scope = data.Get(commonScope).(string)
	return &assignmentRule
}
