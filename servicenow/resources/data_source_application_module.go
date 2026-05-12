package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceApplicationModule reads an Application Module in ServiceNow.
func DataSourceApplicationModule() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceApplicationModule().Schema
	setOnlyRequiredSchema(resourceSchema, applicationModuleTitle)

	return &schema.Resource{
		Description: "`servicenow_application_module` data source can be used to retrieve information of a single application module in ServiceNow by title.",
		ReadContext: readDataSourceApplicationModule,
		Schema:      resourceSchema,
	}
}

func readDataSourceApplicationModule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	applicationModule := &client.ApplicationModule{}
	if err := snowClient.GetObjectByTitle(ctx, client.EndpointApplicationModule, data.Get(applicationModuleTitle).(string), applicationModule); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromApplicationModule(data, applicationModule)

	return nil
}
