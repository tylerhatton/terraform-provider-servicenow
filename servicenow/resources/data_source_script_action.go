package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceScriptAction reads information about a single script action in ServiceNow.
func DataSourceScriptAction() *schema.Resource {
	resourceSchema := ResourceScriptAction().Schema
	setOnlyRequiredSchema(resourceSchema, scriptActionName)

	return &schema.Resource{
		Description: "`servicenow_script_action` data source can be used to retrieve information of a single script action in ServiceNow by name.",
		ReadContext: readDataSourceScriptAction,
		Schema:      resourceSchema,
	}
}

func readDataSourceScriptAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scriptAction := &client.ScriptAction{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointScriptAction, data.Get(scriptActionName).(string), scriptAction); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromScriptAction(data, scriptAction)

	return nil
}
