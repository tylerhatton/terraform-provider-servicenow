package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceRole reads the informations about a single Role in ServiceNow.
func DataSourceRole() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceRole().Schema
	setOnlyRequiredSchema(resourceSchema, roleName)

	return &schema.Resource{
		Description: "`servicenow_role` data source can be used to retrieve information of a single role in ServiceNow by Sys ID",
		ReadContext: readDataSourceRole,
		Schema:      resourceSchema,
	}
}

func readDataSourceRole(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	role := &client.Role{}
	if err := snowClient.GetObjectByName(client.EndpointRole, data.Get(roleName).(string), role); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromRole(data, role)

	return nil
}
