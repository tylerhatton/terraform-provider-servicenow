package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceUIPolicy reads a UI Policy in ServiceNow.
func DataSourceUIPolicy() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceUIPolicy().Schema
	// Look up UI policies by short_description since they don't have a name field.
	setOnlyRequiredSchema(resourceSchema, uiPolicyShortDescription)

	return &schema.Resource{
		Description: "`servicenow_ui_policy` data source can be used to retrieve information of a single UI policy in ServiceNow by its short description.",
		ReadContext: readDataSourceUIPolicy,
		Schema:      resourceSchema,
	}
}

func readDataSourceUIPolicy(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiPolicy := &client.UIPolicy{}
	if err := snowClient.GetObjectByQuery(ctx, client.EndpointUIPolicy, "short_description="+data.Get(uiPolicyShortDescription).(string), uiPolicy); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIPolicy(data, uiPolicy)

	return nil
}
