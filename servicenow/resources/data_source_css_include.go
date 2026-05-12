package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceCSSInclude reads a Service Portal CSS Include in ServiceNow.
func DataSourceCSSInclude() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceCSSInclude().Schema
	setOnlyRequiredSchema(resourceSchema, cssIncludeName)

	return &schema.Resource{
		Description: "`servicenow_css_include` data source can be used to retrieve information of a single Service Portal CSS include in ServiceNow by name.",
		ReadContext: readDataSourceCSSInclude,
		Schema:      resourceSchema,
	}
}

func readDataSourceCSSInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	cssInclude := &client.CSSInclude{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointCSSInclude, data.Get(cssIncludeName).(string), cssInclude); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromCSSInclude(data, cssInclude)

	return nil
}
