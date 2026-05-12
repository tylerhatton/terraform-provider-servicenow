package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceClientScript reads a Client Script in ServiceNow.
func DataSourceClientScript() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceClientScript().Schema
	setOnlyRequiredSchema(resourceSchema, clientScriptName)

	return &schema.Resource{
		Description: "`servicenow_client_script` data source can be used to retrieve information of a single client script in ServiceNow by name.",
		ReadContext: readDataSourceClientScript,
		Schema:      resourceSchema,
	}
}

func readDataSourceClientScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	clientScript := &client.ClientScript{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointClientScript, data.Get(clientScriptName).(string), clientScript); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromClientScript(data, clientScript)

	return nil
}
