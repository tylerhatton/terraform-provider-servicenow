package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceRole reads the informations about a single Role in ServiceNow.
func DataSourceRole() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceRole().Schema
	setOnlyRequiredSchema(resourceSchema, roleName)

	return &schema.Resource{
		Description: "`servicenow_role` data source can be used to retrieve information of a single role in ServiceNow by Sys ID",
		Read:        readDataSourceRole,
		Schema:      resourceSchema,
	}
}

func readDataSourceRole(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	role := &client.Role{}
	if err := snowClient.GetObjectByName(client.EndpointRole, data.Get(roleName).(string), role); err != nil {
		data.SetId("")
		return err
	}

	resourceFromRole(data, role)

	return nil
}
