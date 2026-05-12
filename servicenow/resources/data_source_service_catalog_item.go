package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceServiceCatalogItem reads the information about a single Service Catalog Item in ServiceNow.
func DataSourceServiceCatalogItem() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceServiceCatalogItem().Schema
	setOnlyRequiredSchema(resourceSchema, serviceCatalogItemName)

	return &schema.Resource{
		Description: "`servicenow_service_catalog_item` data source can be used to retrieve information of a single service catalog item in ServiceNow by name.",
		ReadContext: readDataSourceServiceCatalogItem,
		Schema:      resourceSchema,
	}
}

func readDataSourceServiceCatalogItem(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogItem := &client.ServiceCatalogItem{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointServiceCatalogItem, data.Get(serviceCatalogItemName).(string), serviceCatalogItem); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromServiceCatalogItem(data, serviceCatalogItem)

	return nil
}
