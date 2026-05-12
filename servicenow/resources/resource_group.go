package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const groupName = "name"
const groupDescription = "description"
const groupEmail = "email"
const groupManager = "manager"
const groupParent = "parent"
const groupCostCenter = "cost_center"
const groupType = "type"
const groupActive = "active"
const groupDefaultAssignee = "default_assignee"
const groupIncludeMembers = "include_members"
const groupExcludeManager = "exclude_manager"

// ResourceGroup manages a user group in ServiceNow.
func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_group` manages a user group within ServiceNow.",

		CreateContext: createResourceGroup,
		ReadContext:   readResourceGroup,
		UpdateContext: updateResourceGroup,
		DeleteContext: deleteResourceGroup,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			groupName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the user group.",
			},
			groupDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Free-form description of the user group.",
			},
			groupEmail: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Distribution email address associated with the group.",
			},
			groupManager: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the user record that manages this group.",
			},
			groupParent: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the parent group record for hierarchical grouping.",
			},
			groupCostCenter: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the cost center associated with the group.",
			},
			groupType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys IDs of the group type records (sys_user_group_type) associated with this group.",
			},
			groupActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the group is active.",
			},
			groupDefaultAssignee: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the user record that is the default assignee for items routed to this group.",
			},
			groupIncludeMembers: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether nested group members are included when resolving membership.",
			},
			groupExcludeManager: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the group manager is excluded from the member resolution.",
			},
		},
	}
}

func readResourceGroup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	group := &client.Group{}
	if err := snowClient.GetObject(ctx, client.EndpointGroup, data.Id(), group); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromGroup(data, group)

	return nil
}

func createResourceGroup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	group := resourceToGroup(data)
	if err := snowClient.CreateObject(ctx, client.EndpointGroup, group); err != nil {
		return diag.FromErr(err)
	}

	resourceFromGroup(data, group)

	return readResourceGroup(ctx, data, serviceNowClient)
}

func updateResourceGroup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointGroup, resourceToGroup(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceGroup(ctx, data, serviceNowClient)
}

func deleteResourceGroup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointGroup, data.Id()))
}

func resourceFromGroup(data *schema.ResourceData, group *client.Group) {
	data.SetId(group.ID)
	data.Set(groupName, group.Name)
	data.Set(groupDescription, group.Description)
	data.Set(groupEmail, group.Email)
	data.Set(groupManager, group.Manager)
	data.Set(groupParent, group.Parent)
	data.Set(groupCostCenter, group.CostCenter)
	data.Set(groupType, group.Type)
	data.Set(groupActive, group.Active)
	data.Set(groupDefaultAssignee, group.DefaultAssignee)
	data.Set(groupIncludeMembers, group.IncludeMembers)
	data.Set(groupExcludeManager, group.ExcludeManager)
}

func resourceToGroup(data *schema.ResourceData) *client.Group {
	group := client.Group{
		Name:            data.Get(groupName).(string),
		Description:     data.Get(groupDescription).(string),
		Email:           data.Get(groupEmail).(string),
		Manager:         data.Get(groupManager).(string),
		Parent:          data.Get(groupParent).(string),
		CostCenter:      data.Get(groupCostCenter).(string),
		Type:            data.Get(groupType).(string),
		Active:          data.Get(groupActive).(bool),
		DefaultAssignee: data.Get(groupDefaultAssignee).(string),
		IncludeMembers:  data.Get(groupIncludeMembers).(bool),
		ExcludeManager:  data.Get(groupExcludeManager).(bool),
	}
	group.ID = data.Id()
	return &group
}
