package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceScheduledJob reads information about a single scheduled job in ServiceNow.
func DataSourceScheduledJob() *schema.Resource {
	resourceSchema := ResourceScheduledJob().Schema
	setOnlyRequiredSchema(resourceSchema, scheduledJobName)

	return &schema.Resource{
		Description: "`servicenow_scheduled_job` data source can be used to retrieve information of a single scheduled job in ServiceNow by name.",
		ReadContext: readDataSourceScheduledJob,
		Schema:      resourceSchema,
	}
}

func readDataSourceScheduledJob(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scheduledJob := &client.ScheduledJob{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointScheduledJob, data.Get(scheduledJobName).(string), scheduledJob); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromScheduledJob(data, scheduledJob)

	return nil
}
