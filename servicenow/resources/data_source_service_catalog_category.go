package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceServiceCatalogCategory reads a service catalog category in ServiceNow.
func DataSourceServiceCatalogCategory() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceServiceCatalogCategory().Schema
	setOnlyRequiredSchema(resourceSchema, serviceCatalogTitle)

	return &schema.Resource{
		Read:   readDataSourceServiceCatalogCategory,
		Schema: resourceSchema,
	}
}

func readDataSourceServiceCatalogCategory(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogCategory := &client.ServiceCatalogCategory{}
	if err := snowClient.GetObjectByTitle(client.EndpointServiceCatalogCategory, data.Get(serviceCatalogCategoryTitle).(string), serviceCatalogCategory); err != nil {
		data.SetId("")
		return err
	}

	resourceFromServiceCatalogCategory(data, serviceCatalogCategory)

	return nil
}
