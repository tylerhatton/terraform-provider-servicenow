package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceContentCSS reads a Content CSS style sheet in ServiceNow.
func DataSourceContentCSS() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceContentCSS().Schema
	setOnlyRequiredSchema(resourceSchema, contentCSSName)

	return &schema.Resource{
		Description: "`servicenow_content_css` data source can be used to retrieve information of a single content management style sheet in ServiceNow by name.",
		ReadContext: readDataSourceContentCSS,
		Schema:      resourceSchema,
	}
}

func readDataSourceContentCSS(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	contentCSS := &client.ContentCSS{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointContentCSS, data.Get(contentCSSName).(string), contentCSS); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromContentCSS(data, contentCSS)

	return nil
}
