package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceUIPage reads the information about a single UI Page in ServiceNow.
func DataSourceUIPage() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceUIPage().Schema
	setOnlyRequiredSchema(resourceSchema, uiPageName)

	return &schema.Resource{
		Description: "`servicenow_ui_page` data source can be used to retrieve information of a single UI Page in ServiceNow by name.",
		ReadContext: readDataSourceUIPage,
		Schema:      resourceSchema,
	}
}

func readDataSourceUIPage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiPage := &client.UIPage{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointUIPage, data.Get(uiPageName).(string), uiPage); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIPage(data, uiPage)

	return nil
}
