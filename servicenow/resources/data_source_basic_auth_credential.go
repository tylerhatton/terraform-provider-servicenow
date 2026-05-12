package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceBasicAuthCredential reads a Basic Auth Credential in ServiceNow.
func DataSourceBasicAuthCredential() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceBasicAuthCredential().Schema
	setOnlyRequiredSchema(resourceSchema, basicAuthCredentialName)

	return &schema.Resource{
		Description: "`servicenow_basic_auth_credential` data source can be used to retrieve information of a single basic auth credential in ServiceNow by name.",
		ReadContext: readDataSourceBasicAuthCredential,
		Schema:      resourceSchema,
	}
}

func readDataSourceBasicAuthCredential(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	basicAuthCredential := &client.BasicAuthCredential{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointBasicAuthCredential, data.Get(basicAuthCredentialName).(string), basicAuthCredential); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromBasicAuthCredential(data, basicAuthCredential)

	return nil
}
