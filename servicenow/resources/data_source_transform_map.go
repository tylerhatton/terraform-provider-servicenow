package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceTransformMap reads information about a single transform map in ServiceNow.
func DataSourceTransformMap() *schema.Resource {
	resourceSchema := ResourceTransformMap().Schema
	setOnlyRequiredSchema(resourceSchema, transformMapName)

	return &schema.Resource{
		Description: "`servicenow_transform_map` data source can be used to retrieve information of a single transform map in ServiceNow by name.",
		ReadContext: readDataSourceTransformMap,
		Schema:      resourceSchema,
	}
}

func readDataSourceTransformMap(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	transformMap := &client.TransformMap{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointTransformMap, data.Get(transformMapName).(string), transformMap); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromTransformMap(data, transformMap)

	return nil
}
