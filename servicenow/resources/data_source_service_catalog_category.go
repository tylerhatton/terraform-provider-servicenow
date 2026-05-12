package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceServiceCatalogCategory reads a service catalog category in ServiceNow.
func DataSourceServiceCatalogCategory() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceServiceCatalogCategory().Schema
	setOnlyRequiredSchema(resourceSchema, serviceCatalogTitle)

	return &schema.Resource{
		Description: "`servicenow_service_catalog_category` data source can be used to retrieve information of a single service catalog category in ServiceNow by Sys ID",
		ReadContext: readDataSourceServiceCatalogCategory,
		Schema:      resourceSchema,
	}
}

func readDataSourceServiceCatalogCategory(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogCategory := &client.ServiceCatalogCategory{}
	if err := snowClient.GetObjectByTitle(client.EndpointServiceCatalogCategory, data.Get(serviceCatalogCategoryTitle).(string), serviceCatalogCategory); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromServiceCatalogCategory(data, serviceCatalogCategory)

	return nil
}
