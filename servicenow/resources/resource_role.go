package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const roleSuffix = "suffix"
const roleDescription = "description"
const roleElevatedPrivilege = "elevated_privilege"
const roleAssignableBy = "assignable_by"
const roleName = "name"

// ResourceRole manages a Role in ServiceNow.
func ResourceRole() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_role` manages a role within ServiceNow.",

		CreateContext: createResourceRole,
		ReadContext:   readResourceRole,
		UpdateContext: updateResourceRole,
		DeleteContext: deleteResourceRole,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			roleSuffix: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The suffix of the role name. The full role name is prefixed with the application scope.",
			},
			roleDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the role and its purpose.",
			},
			roleElevatedPrivilege: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the role requires elevated privilege to be granted.",
			},
			roleAssignableBy: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-separated list of role names that can assign this role to users.",
			},
			roleName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full name of the role including application scope prefix.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	role := &client.Role{}
	if err := snowClient.GetObject(client.EndpointRole, data.Id(), role); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromRole(data, role)

	return nil
}

func createResourceRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	role := resourceToRole(data)
	if err := snowClient.CreateObject(client.EndpointRole, role); err != nil {
		return diag.FromErr(err)
	}

	resourceFromRole(data, role)

	return readResourceRole(ctx, data, serviceNowClient)
}

func updateResourceRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointRole, resourceToRole(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceRole(ctx, data, serviceNowClient)
}

func deleteResourceRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointRole, data.Id()))
}

func resourceFromRole(data *schema.ResourceData, role *client.Role) {
	data.SetId(role.ID)
	data.Set(roleDescription, role.Description)
	data.Set(roleSuffix, role.Suffix)
	data.Set(roleElevatedPrivilege, role.ElevatedPrivilege)
	data.Set(roleAssignableBy, role.AssignableBy)
	data.Set(roleName, role.Name)
	data.Set(commonProtectionPolicy, role.ProtectionPolicy)
	data.Set(commonScope, role.Scope)
}

func resourceToRole(data *schema.ResourceData) *client.Role {
	role := client.Role{
		Suffix:            data.Get(roleSuffix).(string),
		Description:       data.Get(roleDescription).(string),
		ElevatedPrivilege: data.Get(roleElevatedPrivilege).(bool),
		AssignableBy:      data.Get(roleAssignableBy).(string),
	}
	role.ID = data.Id()
	role.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	role.Scope = data.Get(commonScope).(string)
	return &role
}
