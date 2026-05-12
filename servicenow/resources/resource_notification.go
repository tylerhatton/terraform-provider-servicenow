package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const notificationName = "name"
const notificationTable = "table"
const notificationEventName = "event_name"
const notificationCondition = "condition"
const notificationActive = "active"
const notificationDescription = "description"
const notificationSubject = "subject"
const notificationMessageHTML = "message_html"
const notificationMessagePlain = "message_plain"
const notificationFrom = "from"
const notificationReplyTo = "reply_to"
const notificationRecipientUsers = "recipient_users"
const notificationRecipientGroups = "recipient_groups"
const notificationRecipientRoles = "recipient_roles"
const notificationSendSelf = "send_self"
const notificationIncludeAttachments = "include_attachments"
const notificationCategory = "category"
const notificationWeight = "weight"
const notificationType = "type"
const notificationMandatory = "mandatory"

// ResourceNotification manages an email notification in ServiceNow.
func ResourceNotification() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_notification` manages an email notification (sysevent_email_action) within ServiceNow.",

		CreateContext: createResourceNotification,
		ReadContext:   readResourceNotification,
		UpdateContext: updateResourceNotification,
		DeleteContext: deleteResourceNotification,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			notificationName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the email notification.",
			},
			notificationTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Table that the notification fires on.",
			},
			notificationEventName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of the event that triggers this notification.",
			},
			notificationCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Encoded query condition that must be met for the notification to fire.",
			},
			notificationActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this notification is enabled.",
			},
			notificationDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the notification and its purpose.",
			},
			notificationSubject: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Subject line for the email notification.",
			},
			notificationMessageHTML: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HTML body of the email notification.",
			},
			notificationMessagePlain: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Plain text body of the email notification.",
			},
			notificationFrom: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Email address that the notification will be sent from.",
			},
			notificationReplyTo: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Reply-to email address for the notification.",
			},
			notificationRecipientUsers: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-separated list of user sys_ids that should receive the notification.",
			},
			notificationRecipientGroups: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-separated list of group sys_ids that should receive the notification.",
			},
			notificationRecipientRoles: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-separated list of recipient roles or recipient fields.",
			},
			notificationSendSelf: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, send the notification to the user that triggered the event.",
			},
			notificationIncludeAttachments: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, include attachments from the triggering record on the notification email.",
			},
			notificationCategory: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "sys_id of the notification category record. ServiceNow defaults this from the glide.notification.default_category system property when omitted.",
			},
			notificationWeight: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Numeric weight that controls notification precedence.",
			},
			notificationType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "EMAIL",
				Description: "Notification type (e.g. EMAIL).",
			},
			notificationMandatory: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, recipients cannot unsubscribe from this notification.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceNotification(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	notification := &client.Notification{}
	if err := snowClient.GetObject(ctx, client.EndpointNotification, data.Id(), notification); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromNotification(data, notification)

	return nil
}

func createResourceNotification(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	notification := resourceToNotification(data)
	if err := snowClient.CreateObject(ctx, client.EndpointNotification, notification); err != nil {
		return diag.FromErr(err)
	}

	resourceFromNotification(data, notification)

	return readResourceNotification(ctx, data, serviceNowClient)
}

func updateResourceNotification(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointNotification, resourceToNotification(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceNotification(ctx, data, serviceNowClient)
}

func deleteResourceNotification(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointNotification, data.Id()))
}

func resourceFromNotification(data *schema.ResourceData, notification *client.Notification) {
	data.SetId(notification.ID)
	data.Set(notificationName, notification.Name)
	data.Set(notificationTable, notification.Table)
	data.Set(notificationEventName, notification.EventName)
	data.Set(notificationCondition, notification.Condition)
	data.Set(notificationActive, notification.Active)
	data.Set(notificationDescription, notification.Description)
	data.Set(notificationSubject, notification.Subject)
	data.Set(notificationMessageHTML, notification.MessageHTML)
	data.Set(notificationMessagePlain, notification.MessagePlain)
	data.Set(notificationFrom, notification.From)
	data.Set(notificationReplyTo, notification.ReplyTo)
	data.Set(notificationRecipientUsers, notification.RecipientUsers)
	data.Set(notificationRecipientGroups, notification.RecipientGroups)
	data.Set(notificationRecipientRoles, notification.RecipientRoles)
	data.Set(notificationSendSelf, notification.SendSelf)
	data.Set(notificationIncludeAttachments, notification.IncludeAttachments)
	data.Set(notificationCategory, notification.Category)
	data.Set(notificationWeight, notification.Weight)
	data.Set(notificationType, notification.Type)
	data.Set(notificationMandatory, notification.Mandatory)
	data.Set(commonProtectionPolicy, notification.ProtectionPolicy)
	data.Set(commonScope, notification.Scope)
}

func resourceToNotification(data *schema.ResourceData) *client.Notification {
	notification := client.Notification{
		Name:               data.Get(notificationName).(string),
		Table:              data.Get(notificationTable).(string),
		EventName:          data.Get(notificationEventName).(string),
		Condition:          data.Get(notificationCondition).(string),
		Active:             data.Get(notificationActive).(bool),
		Description:        data.Get(notificationDescription).(string),
		Subject:            data.Get(notificationSubject).(string),
		MessageHTML:        data.Get(notificationMessageHTML).(string),
		MessagePlain:       data.Get(notificationMessagePlain).(string),
		From:               data.Get(notificationFrom).(string),
		ReplyTo:            data.Get(notificationReplyTo).(string),
		RecipientUsers:     data.Get(notificationRecipientUsers).(string),
		RecipientGroups:    data.Get(notificationRecipientGroups).(string),
		RecipientRoles:     data.Get(notificationRecipientRoles).(string),
		SendSelf:           data.Get(notificationSendSelf).(bool),
		IncludeAttachments: data.Get(notificationIncludeAttachments).(bool),
		Category:           data.Get(notificationCategory).(string),
		Weight:             data.Get(notificationWeight).(int),
		Type:               data.Get(notificationType).(string),
		Mandatory:          data.Get(notificationMandatory).(bool),
	}
	notification.ID = data.Id()
	notification.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	notification.Scope = data.Get(commonScope).(string)
	return &notification
}
