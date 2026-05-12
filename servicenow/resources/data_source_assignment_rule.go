package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceAssignmentRule reads information about a single assignment rule in ServiceNow.
func DataSourceAssignmentRule() *schema.Resource {
	resourceSchema := ResourceAssignmentRule().Schema
	setOnlyRequiredSchema(resourceSchema, assignmentRuleName)

	return &schema.Resource{
		Description: "`servicenow_assignment_rule` data source can be used to retrieve information of a single assignment rule in ServiceNow by name.",
		ReadContext: readDataSourceAssignmentRule,
		Schema:      resourceSchema,
	}
}

func readDataSourceAssignmentRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	assignmentRule := &client.AssignmentRule{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointAssignmentRule, data.Get(assignmentRuleName).(string), assignmentRule); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromAssignmentRule(data, assignmentRule)

	return nil
}
