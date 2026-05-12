package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const groupMemberUser = "user"
const groupMemberGroup = "group"

// ResourceGroupMember manages the membership of a user within a group in ServiceNow.
func ResourceGroupMember() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_group_member` manages the membership of a user in a user group within ServiceNow.",

		CreateContext: createResourceGroupMember,
		ReadContext:   readResourceGroupMember,
		DeleteContext: deleteResourceGroupMember,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			groupMemberUser: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sys ID of the user record being added to the group.",
			},
			groupMemberGroup: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sys ID of the user group record receiving the member.",
			},
		},
	}
}

func readResourceGroupMember(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	groupMember := &client.GroupMember{}
	if err := snowClient.GetObject(ctx, client.EndpointGroupMember, data.Id(), groupMember); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromGroupMember(data, groupMember)

	return nil
}

func createResourceGroupMember(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	groupMember := resourceToGroupMember(data)
	if err := snowClient.CreateObject(ctx, client.EndpointGroupMember, groupMember); err != nil {
		return diag.FromErr(err)
	}

	resourceFromGroupMember(data, groupMember)

	return readResourceGroupMember(ctx, data, serviceNowClient)
}

func deleteResourceGroupMember(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointGroupMember, data.Id()))
}

func resourceFromGroupMember(data *schema.ResourceData, groupMember *client.GroupMember) {
	data.SetId(groupMember.ID)
	data.Set(groupMemberUser, groupMember.User)
	data.Set(groupMemberGroup, groupMember.Group)
}

func resourceToGroupMember(data *schema.ResourceData) *client.GroupMember {
	groupMember := client.GroupMember{
		User:  data.Get(groupMemberUser).(string),
		Group: data.Get(groupMemberGroup).(string),
	}
	groupMember.ID = data.Id()
	return &groupMember
}
