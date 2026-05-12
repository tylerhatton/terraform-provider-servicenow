package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceApplication reads an Application in ServiceNow.
func DataSourceApplication() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceApplication().Schema
	setOnlyRequiredSchema(resourceSchema, applicationName)

	return &schema.Resource{
		Description: "`servicenow_application` data source can be used to retrieve information of a single application in ServiceNow by Sys ID",
		ReadContext: readDataSourceApplication,
		Schema:      resourceSchema,
	}
}

func readDataSourceApplication(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	application := &client.Application{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointApplication, data.Get(applicationName).(string), application); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromApplication(data, application)

	return nil
}
