package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceUIScript reads the information about a single UI Script in ServiceNow.
func DataSourceUIScript() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceUIScript().Schema
	setOnlyRequiredSchema(resourceSchema, uiScriptName)

	return &schema.Resource{
		Description: "`servicenow_ui_script` data source can be used to retrieve information of a single UI Script in ServiceNow by name.",
		ReadContext: readDataSourceUIScript,
		Schema:      resourceSchema,
	}
}

func readDataSourceUIScript(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiScript := &client.UIScript{}
	// The UI Script resource's "name" field is mapped to ServiceNow's "script_name" column,
	// which is the user-visible display name. Query by that column rather than the system
	// "name" column (which holds the computed API name).
	if err := snowClient.GetObjectByQuery(ctx, client.EndpointUIScript, "script_name="+data.Get(uiScriptName).(string), uiScript); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIScript(data, uiScript)

	return nil
}
