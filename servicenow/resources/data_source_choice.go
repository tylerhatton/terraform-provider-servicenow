package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceChoice reads a single choice list value in ServiceNow.
func DataSourceChoice() *schema.Resource {
	// Copy the schema from the resource and convert it to a data source schema.
	resourceSchema := ResourceChoice().Schema
	setOnlyRequiredSchema(resourceSchema, choiceName)

	// Name, element, and value together uniquely identify a sys_choice entry.
	resourceSchema[choiceElement].Computed = false
	resourceSchema[choiceElement].Required = true
	resourceSchema[choiceValue].Computed = false
	resourceSchema[choiceValue].Required = true

	return &schema.Resource{
		Description: "`servicenow_choice` data source can be used to retrieve information of a single choice list value in ServiceNow by the combination of table name, field name and value.",
		ReadContext: readDataSourceChoice,
		Schema:      resourceSchema,
	}
}

func readDataSourceChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	choice := &client.Choice{}
	name := data.Get(choiceName).(string)
	element := data.Get(choiceElement).(string)
	value := data.Get(choiceValue).(string)

	query := "name=" + name + "^element=" + element + "^value=" + value
	if err := snowClient.GetObjectByQuery(ctx, client.EndpointChoice, query, choice); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromChoice(data, choice)

	return nil
}
