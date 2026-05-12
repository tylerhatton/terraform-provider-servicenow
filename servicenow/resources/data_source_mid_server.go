package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceMidServer reads a MID server (ECC agent) record in ServiceNow.
func DataSourceMidServer() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceMidServer().Schema
	setOnlyRequiredSchema(resourceSchema, midServerName)

	return &schema.Resource{
		Description: "`servicenow_mid_server` data source can be used to retrieve information of a single MID server (ECC agent) in ServiceNow by name.",
		ReadContext: readDataSourceMidServer,
		Schema:      resourceSchema,
	}
}

func readDataSourceMidServer(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	midServer := &client.MidServer{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointMidServer, data.Get(midServerName).(string), midServer); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromMidServer(data, midServer)

	return nil
}
