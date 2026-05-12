package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceAlias reads a connection or credential Alias in ServiceNow.
func DataSourceAlias() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceAlias().Schema
	setOnlyRequiredSchema(resourceSchema, aliasName)

	return &schema.Resource{
		Description: "`servicenow_alias` data source can be used to retrieve information of a single connection or credential alias in ServiceNow by name.",
		ReadContext: readDataSourceAlias,
		Schema:      resourceSchema,
	}
}

func readDataSourceAlias(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	alias := &client.Alias{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointAlias, data.Get(aliasName).(string), alias); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromAlias(data, alias)

	return nil
}
