package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceFlow reads a Flow Designer flow record in ServiceNow.
func DataSourceFlow() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceFlow().Schema
	setOnlyRequiredSchema(resourceSchema, flowName)

	return &schema.Resource{
		Description: "`servicenow_flow` data source can be used to retrieve information of a single Flow Designer flow in ServiceNow by name.",
		ReadContext: readDataSourceFlow,
		Schema:      resourceSchema,
	}
}

func readDataSourceFlow(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	flow := &client.Flow{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointFlow, data.Get(flowName).(string), flow); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromFlow(data, flow)

	return nil
}
