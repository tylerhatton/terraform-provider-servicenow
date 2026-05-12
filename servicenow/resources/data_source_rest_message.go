package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceRestMessage reads a REST Message configuration in ServiceNow.
func DataSourceRestMessage() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceRestMessage().Schema
	setOnlyRequiredSchema(resourceSchema, restMessageName)

	return &schema.Resource{
		Description: "`servicenow_rest_message` data source can be used to retrieve information of a single REST message in ServiceNow by name.",
		ReadContext: readDataSourceRestMessage,
		Schema:      resourceSchema,
	}
}

func readDataSourceRestMessage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMessage := &client.RestMessage{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointRestMessage, data.Get(restMessageName).(string), restMessage); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromRestMessage(data, restMessage)

	return nil
}
