package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceWidget reads the information about a single Widget in ServiceNow.
func DataSourceWidget() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceWidget().Schema
	setOnlyRequiredSchema(resourceSchema, widgetName)

	return &schema.Resource{
		Description: "`servicenow_widget` data source can be used to retrieve information of a single Widget in ServiceNow by name.",
		ReadContext: readDataSourceWidget,
		Schema:      resourceSchema,
	}
}

func readDataSourceWidget(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	widget := &client.Widget{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointWidget, data.Get(widgetName).(string), widget); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromWidget(data, widget)

	return nil
}
