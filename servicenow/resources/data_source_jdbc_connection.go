package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceJdbcConnection reads a JDBC connection configuration in ServiceNow.
func DataSourceJdbcConnection() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceJdbcConnection().Schema
	setOnlyRequiredSchema(resourceSchema, jdbcConnectionName)

	return &schema.Resource{
		Description: "`servicenow_jdbc_connection` data source can be used to retrieve information of a single JDBC connection in ServiceNow by name.",
		ReadContext: readDataSourceJdbcConnection,
		Schema:      resourceSchema,
	}
}

func readDataSourceJdbcConnection(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	jdbcConnection := &client.JdbcConnection{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointJdbcConnection, data.Get(jdbcConnectionName).(string), jdbcConnection); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromJdbcConnection(data, jdbcConnection)

	return nil
}
