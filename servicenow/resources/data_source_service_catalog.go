package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceServiceCatalog reads the informations about a single ServiceCatalog in ServiceNow.
func DataSourceServiceCatalog() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceServiceCatalog().Schema
	setOnlyRequiredSchema(resourceSchema, serviceCatalogTitle)

	return &schema.Resource{
		Description: "`servicenow_service_catalog` data source can be used to retrieve information of a single service catalog in ServiceNow by Sys ID",
		Read:        readDataSourceServiceCatalog,
		Schema:      resourceSchema,
	}
}

func readDataSourceServiceCatalog(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalog := &client.ServiceCatalog{}
	if err := snowClient.GetObjectByTitle(client.EndpointServiceCatalog, data.Get(serviceCatalogTitle).(string), serviceCatalog); err != nil {
		data.SetId("")
		return err
	}

	resourceFromServiceCatalog(data, serviceCatalog)

	return nil
}
