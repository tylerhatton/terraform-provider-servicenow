package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceScriptInclude reads the information about a single Script Include in ServiceNow.
func DataSourceScriptInclude() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceScriptInclude().Schema
	setOnlyRequiredSchema(resourceSchema, scriptIncludeName)

	return &schema.Resource{
		Description: "`servicenow_script_include` data source can be used to retrieve information of a single script include in ServiceNow by name.",
		ReadContext: readDataSourceScriptInclude,
		Schema:      resourceSchema,
	}
}

func readDataSourceScriptInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scriptInclude := &client.ScriptInclude{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointScriptInclude, data.Get(scriptIncludeName).(string), scriptInclude); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromScriptInclude(data, scriptInclude)

	return nil
}
