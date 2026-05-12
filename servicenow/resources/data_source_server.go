package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceServer reads the information about a single Server CMDB entry in ServiceNow.
func DataSourceServer() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceServer().Schema
	setOnlyRequiredSchema(resourceSchema, serverName)

	return &schema.Resource{
		Description: "`servicenow_server` data source can be used to retrieve information of a single server CMDB entry in ServiceNow by name.",
		ReadContext: readDataSourceServer,
		Schema:      resourceSchema,
	}
}

func readDataSourceServer(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	server := &client.Server{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointServer, data.Get(serverName).(string), server); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromServer(data, server)

	return nil
}
