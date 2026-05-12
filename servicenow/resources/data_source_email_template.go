package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// DataSourceEmailTemplate reads information about a single email template in ServiceNow.
func DataSourceEmailTemplate() *schema.Resource {
	resourceSchema := ResourceEmailTemplate().Schema
	setOnlyRequiredSchema(resourceSchema, emailTemplateName)

	return &schema.Resource{
		Description: "`servicenow_email_template` data source can be used to retrieve information of a single email template in ServiceNow by name.",
		ReadContext: readDataSourceEmailTemplate,
		Schema:      resourceSchema,
	}
}

func readDataSourceEmailTemplate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	emailTemplate := &client.EmailTemplate{}
	if err := snowClient.GetObjectByName(ctx, client.EndpointEmailTemplate, data.Get(emailTemplateName).(string), emailTemplate); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromEmailTemplate(data, emailTemplate)

	return nil
}
