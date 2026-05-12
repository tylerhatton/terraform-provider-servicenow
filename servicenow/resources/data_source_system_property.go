package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceSystemProperty reads the informations about a single SystemProperty in ServiceNow.
func DataSourceSystemProperty() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceSystemProperty().Schema
	// Look up system properties by suffix since the name field is not populated by the JSONv2 API.
	setOnlyRequiredSchema(resourceSchema, systemPropertySuffix)

	return &schema.Resource{
		Description: "`servicenow_system_property` data source can be used to retrieve information of a single system property in ServiceNow by its suffix",
		ReadContext: readDataSourceSystemProperty,
		Schema:      resourceSchema,
	}
}

func readDataSourceSystemProperty(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	systemProperty := &client.SystemProperty{}
	if err := snowClient.GetObjectByQuery(client.EndpointSystemProperty, "suffix="+data.Get(systemPropertySuffix).(string), systemProperty); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromSystemProperty(data, systemProperty)

	return nil
}
