package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceHttpConnection reads an HTTP Connection configuration in ServiceNow.
func DataSourceHttpConnection() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceHttpConnection().Schema
	setOnlyRequiredSchema(resourceSchema, httpConnectionName)

	return &schema.Resource{
		Description: "`servicenow_http_connection` data source can be used to retrieve information of a single HTTP connection in ServiceNow by name.",
		ReadContext: readDataSourceHttpConnection,
		Schema:      resourceSchema,
	}
}

func readDataSourceHttpConnection(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	httpConnection := &client.HttpConnection{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointHttpConnection, data.Get(httpConnectionName).(string), httpConnection); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromHttpConnection(data, httpConnection)

	return nil
}
