package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceJsInclude reads a JS Include in ServiceNow.
func DataSourceJsInclude() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceJsInclude().Schema
	setOnlyRequiredSchema(resourceSchema, jsIncludeDisplayName)

	return &schema.Resource{
		Description: "`servicenow_js_include` data source can be used to retrieve information of a single JS include in ServiceNow by display name.",
		ReadContext: readDataSourceJsInclude,
		Schema:      resourceSchema,
	}
}

func readDataSourceJsInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	jsInclude := &client.JsInclude{}
	if err := snowClient.GetObjectByQuery(ctx, client.EndpointJsInclude, "display_name="+data.Get(jsIncludeDisplayName).(string), jsInclude); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromJsInclude(data, jsInclude)

	return nil
}
