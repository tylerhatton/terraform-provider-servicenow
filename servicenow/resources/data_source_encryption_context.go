package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceEncryptionContext reads an encryption context record in ServiceNow.
func DataSourceEncryptionContext() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceEncryptionContext().Schema
	setOnlyRequiredSchema(resourceSchema, encryptionContextName)

	return &schema.Resource{
		Description: "`servicenow_encryption_context` data source can be used to retrieve information of a single encryption context record in ServiceNow by name. Requires the Edge Encryption plugin to be active.",
		ReadContext: readDataSourceEncryptionContext,
		Schema:      resourceSchema,
	}
}

func readDataSourceEncryptionContext(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	encryptionContext := &client.EncryptionContext{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointEncryptionContext, data.Get(encryptionContextName).(string), encryptionContext); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromEncryptionContext(data, encryptionContext)

	return nil
}
