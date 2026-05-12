package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceUser reads the information about a single User in ServiceNow by user_name.
func DataSourceUser() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceUser().Schema
	// Look up users by user_name (the unique login identifier).
	setOnlyRequiredSchema(resourceSchema, userUserName)

	return &schema.Resource{
		Description: "`servicenow_user` data source can be used to retrieve information about a single user in ServiceNow by its user_name.",
		ReadContext: readDataSourceUser,
		Schema:      resourceSchema,
	}
}

func readDataSourceUser(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	user := &client.User{}
	if err := snowClient.GetObjectByQuery(ctx, client.EndpointUser, "user_name="+data.Get(userUserName).(string), user); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUser(data, user)

	return nil
}
