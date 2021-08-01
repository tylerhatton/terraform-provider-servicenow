package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceApplication reads an Application in ServiceNow.
func DataSourceApplication() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceApplication().Schema
	setOnlyRequiredSchema(resourceSchema, applicationName)

	return &schema.Resource{
		Description: "`servicenow_application` data source can be used to retrieve information of a single application in ServiceNow by Sys ID",
		Read:        readDataSourceApplication,
		Schema:      resourceSchema,
	}
}

func readDataSourceApplication(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	application := &client.Application{}
	if err := snowClient.GetObjectByName(client.EndpointApplication, data.Get(applicationName).(string), application); err != nil {
		data.SetId("")
		return err
	}

	resourceFromApplication(data, application)

	return nil
}
