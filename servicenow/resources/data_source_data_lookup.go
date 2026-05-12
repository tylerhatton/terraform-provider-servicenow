package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceDataLookup reads information about a single data lookup definition in ServiceNow.
func DataSourceDataLookup() *schema.Resource {
	resourceSchema := ResourceDataLookup().Schema
	setOnlyRequiredSchema(resourceSchema, dataLookupName)

	return &schema.Resource{
		Description: "`servicenow_data_lookup` data source can be used to retrieve information of a single data lookup definition in ServiceNow by name.",
		ReadContext: readDataSourceDataLookup,
		Schema:      resourceSchema,
	}
}

func readDataSourceDataLookup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	dataLookup := &client.DataLookup{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointDataLookup, data.Get(dataLookupName).(string), dataLookup); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromDataLookup(data, dataLookup)

	return nil
}
