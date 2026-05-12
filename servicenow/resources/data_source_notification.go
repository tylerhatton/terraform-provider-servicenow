package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceNotification reads information about a single email notification in ServiceNow.
func DataSourceNotification() *schema.Resource {
	resourceSchema := ResourceNotification().Schema
	setOnlyRequiredSchema(resourceSchema, notificationName)

	return &schema.Resource{
		Description: "`servicenow_notification` data source can be used to retrieve information of a single email notification in ServiceNow by name.",
		ReadContext: readDataSourceNotification,
		Schema:      resourceSchema,
	}
}

func readDataSourceNotification(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	notification := &client.Notification{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointNotification, data.Get(notificationName).(string), notification); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromNotification(data, notification)

	return nil
}
