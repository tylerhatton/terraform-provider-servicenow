package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const userRoleUser = "user"
const userRoleRole = "role"
const userRoleInherited = "inherited"
const userRoleGrantedBy = "granted_by"
const userRoleState = "state"

// ResourceUserRole manages the association between a user and a role in ServiceNow.
func ResourceUserRole() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_user_role` assigns a role to a user within ServiceNow.",

		CreateContext: createResourceUserRole,
		ReadContext:   readResourceUserRole,
		UpdateContext: updateResourceUserRole,
		DeleteContext: deleteResourceUserRole,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			userRoleUser: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sys ID of the user record receiving the role assignment.",
			},
			userRoleRole: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sys ID of the role record being assigned.",
			},
			userRoleInherited: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether this role was inherited from a group membership rather than directly assigned.",
			},
			userRoleGrantedBy: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the record (such as a group) that granted this role to the user.",
			},
			userRoleState: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "active",
				Description: "Lifecycle state of the role assignment (typically 'active').",
			},
		},
	}
}

func readResourceUserRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	userRole := &client.UserRole{}
	if err := snowClient.GetObject(ctx, client.EndpointUserRole, data.Id(), userRole); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUserRole(data, userRole)

	return nil
}

func createResourceUserRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	userRole := resourceToUserRole(data)
	if err := snowClient.CreateObject(ctx, client.EndpointUserRole, userRole); err != nil {
		return diag.FromErr(err)
	}

	resourceFromUserRole(data, userRole)

	return readResourceUserRole(ctx, data, serviceNowClient)
}

func updateResourceUserRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointUserRole, resourceToUserRole(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceUserRole(ctx, data, serviceNowClient)
}

func deleteResourceUserRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointUserRole, data.Id()))
}

func resourceFromUserRole(data *schema.ResourceData, userRole *client.UserRole) {
	data.SetId(userRole.ID)
	data.Set(userRoleUser, userRole.User)
	data.Set(userRoleRole, userRole.Role)
	data.Set(userRoleInherited, userRole.Inherited)
	data.Set(userRoleGrantedBy, userRole.GrantedBy)
	data.Set(userRoleState, userRole.State)
}

func resourceToUserRole(data *schema.ResourceData) *client.UserRole {
	userRole := client.UserRole{
		User:      data.Get(userRoleUser).(string),
		Role:      data.Get(userRoleRole).(string),
		Inherited: data.Get(userRoleInherited).(bool),
		GrantedBy: data.Get(userRoleGrantedBy).(string),
		State:     data.Get(userRoleState).(string),
	}
	userRole.ID = data.Id()
	return &userRole
}
