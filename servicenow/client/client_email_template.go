package client

// EndpointEmailTemplate is the endpoint to manage email template (sysevent_email_template) records.
const EndpointEmailTemplate = "sysevent_email_template.do"

// EmailTemplate is the json response for an email template in ServiceNow.
//
// Note: sysevent_email_template does not expose description/active/default/event_name
// columns on the base record (some are surfaced from the related sysevent_email_action
// notifications). Only the fields actually backed by the table are modelled here.
type EmailTemplate struct {
	BaseResult
	Name         string `json:"name"`
	Table        string `json:"collection"`
	Subject      string `json:"subject"`
	MessageHTML  string `json:"message_html"`
	MessagePlain string `json:"message"`
	SMSAlternate string `json:"sms_alternate"`
}
