package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceGroup reads the information about a single user group in ServiceNow by name.
func DataSourceGroup() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceGroup().Schema
	// Look up groups by name.
	setOnlyRequiredSchema(resourceSchema, groupName)

	return &schema.Resource{
		Description: "`servicenow_group` data source can be used to retrieve information about a single user group in ServiceNow by its name.",
		ReadContext: readDataSourceGroup,
		Schema:      resourceSchema,
	}
}

func readDataSourceGroup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	group := &client.Group{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointGroup, data.Get(groupName).(string), group); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromGroup(data, group)

	return nil
}
