package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceScriptedRestApi reads the information about a single Scripted REST API in ServiceNow.
func DataSourceScriptedRestApi() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceScriptedRestApi().Schema
	setOnlyRequiredSchema(resourceSchema, scriptedRestApiName)

	return &schema.Resource{
		Description: "`servicenow_scripted_rest_api` data source can be used to retrieve information of a single scripted REST API in ServiceNow by name.",
		ReadContext: readDataSourceScriptedRestApi,
		Schema:      resourceSchema,
	}
}

func readDataSourceScriptedRestApi(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scriptedRestApi := &client.ScriptedRestApi{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointScriptedRestApi, data.Get(scriptedRestApiName).(string), scriptedRestApi); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromScriptedRestApi(data, scriptedRestApi)

	return nil
}
