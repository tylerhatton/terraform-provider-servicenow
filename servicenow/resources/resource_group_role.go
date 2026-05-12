package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const groupRoleGroup = "group"
const groupRoleRole = "role"
const groupRoleInherits = "inherits"

// ResourceGroupRole manages the association between a group and a role in ServiceNow.
func ResourceGroupRole() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_group_role` assigns a role to a user group within ServiceNow.",

		CreateContext: createResourceGroupRole,
		ReadContext:   readResourceGroupRole,
		UpdateContext: updateResourceGroupRole,
		DeleteContext: deleteResourceGroupRole,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			groupRoleGroup: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sys ID of the user group record receiving the role assignment.",
			},
			groupRoleRole: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sys ID of the role record being assigned to the group.",
			},
			groupRoleInherits: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether members of the group inherit this role assignment.",
			},
		},
	}
}

func readResourceGroupRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	groupRole := &client.GroupRole{}
	if err := snowClient.GetObject(ctx, client.EndpointGroupRole, data.Id(), groupRole); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromGroupRole(data, groupRole)

	return nil
}

func createResourceGroupRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	groupRole := resourceToGroupRole(data)
	if err := snowClient.CreateObject(ctx, client.EndpointGroupRole, groupRole); err != nil {
		return diag.FromErr(err)
	}

	resourceFromGroupRole(data, groupRole)

	return readResourceGroupRole(ctx, data, serviceNowClient)
}

func updateResourceGroupRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointGroupRole, resourceToGroupRole(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceGroupRole(ctx, data, serviceNowClient)
}

func deleteResourceGroupRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointGroupRole, data.Id()))
}

func resourceFromGroupRole(data *schema.ResourceData, groupRole *client.GroupRole) {
	data.SetId(groupRole.ID)
	data.Set(groupRoleGroup, groupRole.Group)
	data.Set(groupRoleRole, groupRole.Role)
	data.Set(groupRoleInherits, groupRole.Inherits)
}

func resourceToGroupRole(data *schema.ResourceData) *client.GroupRole {
	groupRole := client.GroupRole{
		Group:    data.Get(groupRoleGroup).(string),
		Role:     data.Get(groupRoleRole).(string),
		Inherits: data.Get(groupRoleInherits).(bool),
	}
	groupRole.ID = data.Id()
	return &groupRole
}
