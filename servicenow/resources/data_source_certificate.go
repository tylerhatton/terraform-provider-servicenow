package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceCertificate reads a certificate record in ServiceNow.
func DataSourceCertificate() *schema.Resource {
	// Copy the schema from the resource.
	resourceSchema := ResourceCertificate().Schema
	setOnlyRequiredSchema(resourceSchema, certificateName)

	return &schema.Resource{
		Description: "`servicenow_certificate` data source can be used to retrieve information of a single certificate record in ServiceNow by name.",
		ReadContext: readDataSourceCertificate,
		Schema:      resourceSchema,
	}
}

func readDataSourceCertificate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	certificate := &client.Certificate{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointCertificate, data.Get(certificateName).(string), certificate); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromCertificate(data, certificate)

	return nil
}
