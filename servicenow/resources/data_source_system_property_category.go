package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceSystemPropertyCategory reads the informations about a single SystemPropertyCategory in ServiceNow.
func DataSourceSystemPropertyCategory() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceSystemPropertyCategory().Schema
	setOnlyRequiredSchema(resourceSchema, systemPropertyCategoryName)

	return &schema.Resource{
		Description: "`servicenow_system_property_category` data source can be used to retrieve information of a single system property category in ServiceNow by Sys ID",
		ReadContext: readDataSourceSystemPropertyCategory,
		Schema:      resourceSchema,
	}
}

func readDataSourceSystemPropertyCategory(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	systemPropertyCategory := &client.SystemPropertyCategory{}
	if err := snowClient.GetObjectByName(client.EndpointSystemPropertyCategory, data.Get(systemPropertyCategoryName).(string), systemPropertyCategory); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromSystemPropertyCategory(data, systemPropertyCategory)

	return nil
}
