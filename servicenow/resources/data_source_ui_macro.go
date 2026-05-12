package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceUIMacro reads the information about a single UI Macro in ServiceNow.
func DataSourceUIMacro() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceUIMacro().Schema
	setOnlyRequiredSchema(resourceSchema, uiMacroName)

	return &schema.Resource{
		Description: "`servicenow_ui_macro` data source can be used to retrieve information of a single UI Macro in ServiceNow by name.",
		ReadContext: readDataSourceUIMacro,
		Schema:      resourceSchema,
	}
}

func readDataSourceUIMacro(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiMacro := &client.UIMacro{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointUIMacro, data.Get(uiMacroName).(string), uiMacro); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIMacro(data, uiMacro)

	return nil
}
