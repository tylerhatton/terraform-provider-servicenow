package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceExtensionPoint reads a scripted Extension Point in ServiceNow.
func DataSourceExtensionPoint() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceExtensionPoint().Schema
	setOnlyRequiredSchema(resourceSchema, extensionPointName)

	return &schema.Resource{
		Description: "`servicenow_extension_point` data source can be used to retrieve information of a single scripted extension point in ServiceNow by name.",
		ReadContext: readDataSourceExtensionPoint,
		Schema:      resourceSchema,
	}
}

func readDataSourceExtensionPoint(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	extensionPoint := &client.ExtensionPoint{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointExtensionPoint, data.Get(extensionPointName).(string), extensionPoint); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromExtensionPoint(data, extensionPoint)

	return nil
}
