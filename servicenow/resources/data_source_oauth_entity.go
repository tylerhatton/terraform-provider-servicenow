package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceOAuthEntity reads an OAuth application entity in ServiceNow.
func DataSourceOAuthEntity() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceOAuthEntity().Schema
	setOnlyRequiredSchema(resourceSchema, oauthEntityName)

	return &schema.Resource{
		Description: "`servicenow_oauth_entity` data source can be used to retrieve information of a single OAuth application entity in ServiceNow by name.",
		ReadContext: readDataSourceOAuthEntity,
		Schema:      resourceSchema,
	}
}

func readDataSourceOAuthEntity(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	oauthEntity := &client.OAuthEntity{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointOAuthEntity, data.Get(oauthEntityName).(string), oauthEntity); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromOAuthEntity(data, oauthEntity)

	return nil
}
