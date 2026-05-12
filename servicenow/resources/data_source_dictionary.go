package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceDictionary reads a single dictionary entry (column definition) in ServiceNow.
func DataSourceDictionary() *schema.Resource {
	// Copy the schema from the resource and convert it to a data source schema.
	resourceSchema := ResourceDictionary().Schema
	setOnlyRequiredSchema(resourceSchema, dictionaryName)

	// Both name and element are required to uniquely identify a dictionary entry.
	resourceSchema[dictionaryElement].Computed = false
	resourceSchema[dictionaryElement].Required = true

	return &schema.Resource{
		Description: "`servicenow_dictionary` data source can be used to retrieve information of a single column definition in ServiceNow by the combination of table name and field name.",
		ReadContext: readDataSourceDictionary,
		Schema:      resourceSchema,
	}
}

func readDataSourceDictionary(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	dictionary := &client.Dictionary{}
	name := data.Get(dictionaryName).(string)
	element := data.Get(dictionaryElement).(string)

	query := "name=" + name + "^element=" + element
	if err := snowClient.GetObjectByQuery(ctx, client.EndpointDictionary, query, dictionary); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromDictionary(data, dictionary)

	return nil
}
