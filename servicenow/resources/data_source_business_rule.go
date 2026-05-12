package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceBusinessRule reads a Business Rule in ServiceNow.
func DataSourceBusinessRule() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceBusinessRule().Schema
	setOnlyRequiredSchema(resourceSchema, businessRuleName)

	return &schema.Resource{
		Description: "`servicenow_business_rule` data source can be used to retrieve information of a single business rule in ServiceNow by name.",
		ReadContext: readDataSourceBusinessRule,
		Schema:      resourceSchema,
	}
}

func readDataSourceBusinessRule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	businessRule := &client.BusinessRule{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointBusinessRule, data.Get(businessRuleName).(string), businessRule); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromBusinessRule(data, businessRule)

	return nil
}
