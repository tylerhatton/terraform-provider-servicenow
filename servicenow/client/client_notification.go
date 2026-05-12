package client

// EndpointNotification is the endpoint to manage email notification (sysevent_email_action) records.
const EndpointNotification = "sysevent_email_action.do"

// Notification is the json response for an email notification in ServiceNow.
type Notification struct {
	BaseResult
	Name               string `json:"name"`
	Table              string `json:"collection"`
	EventName          string `json:"event_name"`
	Condition          string `json:"condition"`
	Active             bool   `json:"active,string"`
	Description        string `json:"description"`
	Subject            string `json:"subject"`
	MessageHTML        string `json:"message_html"`
	MessagePlain       string `json:"message_text"`
	From               string `json:"from"`
	ReplyTo            string `json:"reply_to"`
	RecipientUsers     string `json:"recipient_users"`
	RecipientGroups    string `json:"recipient_groups"`
	RecipientRoles     string `json:"recipient_fields"`
	SendSelf           bool   `json:"send_self,string"`
	IncludeAttachments bool   `json:"include_attachments,string"`
	Category           string `json:"category"`
	Weight             int    `json:"weight,string"`
	Type               string `json:"type"`
	Mandatory          bool   `json:"mandatory,string"`
}
