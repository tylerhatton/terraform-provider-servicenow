package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceSystemProperty reads the informations about a single SystemProperty in ServiceNow.
func DataSourceSystemProperty() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceSystemProperty().Schema
	setOnlyRequiredSchema(resourceSchema, systemPropertyName)

	return &schema.Resource{
		Description: "`servicenow_system_property` data source can be used to retrieve information of a single system property in ServiceNow by Sys ID",
		Read:        readDataSourceSystemProperty,
		Schema:      resourceSchema,
	}
}

func readDataSourceSystemProperty(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	systemProperty := &client.SystemProperty{}
	if err := snowClient.GetObjectByName(client.EndpointSystemProperty, data.Get(systemPropertyName).(string), systemProperty); err != nil {
		data.SetId("")
		return err
	}

	resourceFromSystemProperty(data, systemProperty)

	return nil
}
