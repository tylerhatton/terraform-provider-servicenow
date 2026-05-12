package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceUIAction reads a UI Action in ServiceNow.
func DataSourceUIAction() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceUIAction().Schema
	setOnlyRequiredSchema(resourceSchema, uiActionName)

	return &schema.Resource{
		Description: "`servicenow_ui_action` data source can be used to retrieve information of a single UI action in ServiceNow by name.",
		ReadContext: readDataSourceUIAction,
		Schema:      resourceSchema,
	}
}

func readDataSourceUIAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiAction := &client.UIAction{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointUIAction, data.Get(uiActionName).(string), uiAction); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIAction(data, uiAction)

	return nil
}
