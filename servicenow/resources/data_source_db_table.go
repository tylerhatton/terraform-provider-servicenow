package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceDBTable reads a DB Table in ServiceNow.
func DataSourceDBTable() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceDBTable().Schema
	setOnlyRequiredSchema(resourceSchema, dbTableName)

	return &schema.Resource{
		Description: "`servicenow_db_table` data source can be used to retrieve information of a single database in ServiceNow by Name",
		ReadContext: readDataSourceDBTable,
		Schema:      resourceSchema,
	}
}

func readDataSourceDBTable(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	dbTable := &client.DBTable{}
	if err := snowClient.GetObjectByName(client.EndpointDBTable, data.Get(dbTableName).(string), dbTable); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromDBTable(data, dbTable)

	return nil
}
