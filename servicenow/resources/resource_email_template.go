package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const emailTemplateName = "name"
const emailTemplateTable = "table"
const emailTemplateSubject = "subject"
const emailTemplateMessageHTML = "message_html"
const emailTemplateMessagePlain = "message_plain"
const emailTemplateSMSAlternate = "sms_alternate"

// ResourceEmailTemplate manages an email template in ServiceNow.
func ResourceEmailTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_email_template` manages an email template (sysevent_email_template) within ServiceNow.",

		CreateContext: createResourceEmailTemplate,
		ReadContext:   readResourceEmailTemplate,
		UpdateContext: updateResourceEmailTemplate,
		DeleteContext: deleteResourceEmailTemplate,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			emailTemplateName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the email template.",
			},
			emailTemplateTable: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Table that the email template is associated with (maps to ServiceNow's `collection` column).",
			},
			emailTemplateSubject: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Subject line for the email template.",
			},
			emailTemplateMessageHTML: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HTML body of the email template.",
			},
			emailTemplateMessagePlain: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Plain text body of the email template (maps to ServiceNow's `message` column).",
			},
			emailTemplateSMSAlternate: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "SMS alternate text for the email template.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceEmailTemplate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	emailTemplate := &client.EmailTemplate{}
	if err := snowClient.GetObject(ctx, client.EndpointEmailTemplate, data.Id(), emailTemplate); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromEmailTemplate(data, emailTemplate)

	return nil
}

func createResourceEmailTemplate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	emailTemplate := resourceToEmailTemplate(data)
	if err := snowClient.CreateObject(ctx, client.EndpointEmailTemplate, emailTemplate); err != nil {
		return diag.FromErr(err)
	}

	resourceFromEmailTemplate(data, emailTemplate)

	return readResourceEmailTemplate(ctx, data, serviceNowClient)
}

func updateResourceEmailTemplate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointEmailTemplate, resourceToEmailTemplate(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceEmailTemplate(ctx, data, serviceNowClient)
}

func deleteResourceEmailTemplate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointEmailTemplate, data.Id()))
}

func resourceFromEmailTemplate(data *schema.ResourceData, emailTemplate *client.EmailTemplate) {
	data.SetId(emailTemplate.ID)
	data.Set(emailTemplateName, emailTemplate.Name)
	data.Set(emailTemplateTable, emailTemplate.Table)
	data.Set(emailTemplateSubject, emailTemplate.Subject)
	data.Set(emailTemplateMessageHTML, emailTemplate.MessageHTML)
	data.Set(emailTemplateMessagePlain, emailTemplate.MessagePlain)
	data.Set(emailTemplateSMSAlternate, emailTemplate.SMSAlternate)
	data.Set(commonProtectionPolicy, emailTemplate.ProtectionPolicy)
	data.Set(commonScope, emailTemplate.Scope)
}

func resourceToEmailTemplate(data *schema.ResourceData) *client.EmailTemplate {
	emailTemplate := client.EmailTemplate{
		Name:         data.Get(emailTemplateName).(string),
		Table:        data.Get(emailTemplateTable).(string),
		Subject:      data.Get(emailTemplateSubject).(string),
		MessageHTML:  data.Get(emailTemplateMessageHTML).(string),
		MessagePlain: data.Get(emailTemplateMessagePlain).(string),
		SMSAlternate: data.Get(emailTemplateSMSAlternate).(string),
	}
	emailTemplate.ID = data.Id()
	emailTemplate.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	emailTemplate.Scope = data.Get(commonScope).(string)
	return &emailTemplate
}
