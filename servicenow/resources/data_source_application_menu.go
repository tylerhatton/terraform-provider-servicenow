package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceApplicationMenu reads an Application Menu in ServiceNow.
func DataSourceApplicationMenu() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceApplicationMenu().Schema
	setOnlyRequiredSchema(resourceSchema, applicationMenuTitle)

	return &schema.Resource{
		Description: "`servicenow_application_menu` data source can be used to retrieve information of a single application menu in ServiceNow by title.",
		ReadContext: readDataSourceApplicationMenu,
		Schema:      resourceSchema,
	}
}

func readDataSourceApplicationMenu(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	applicationMenu := &client.ApplicationMenu{}
	if err := snowClient.GetObjectByTitle(ctx, client.EndpointApplicationMenu, data.Get(applicationMenuTitle).(string), applicationMenu); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromApplicationMenu(data, applicationMenu)

	return nil
}
