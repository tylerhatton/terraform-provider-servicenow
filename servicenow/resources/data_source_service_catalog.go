package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceServiceCatalog reads the informations about a single ServiceCatalog in ServiceNow.
func DataSourceServiceCatalog() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceServiceCatalog().Schema
	setOnlyRequiredSchema(resourceSchema, serviceCatalogTitle)

	return &schema.Resource{
		Description: "`servicenow_service_catalog` data source can be used to retrieve information of a single service catalog in ServiceNow by Sys ID",
		ReadContext: readDataSourceServiceCatalog,
		Schema:      resourceSchema,
	}
}

func readDataSourceServiceCatalog(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalog := &client.ServiceCatalog{}
	if err := snowClient.GetObjectByTitle(client.EndpointServiceCatalog, data.Get(serviceCatalogTitle).(string), serviceCatalog); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromServiceCatalog(data, serviceCatalog)

	return nil
}
