package client

// EndpointUIAction is the endpoint to manage UI action records.
const EndpointUIAction = "sys_ui_action.do"

// UIAction is the json response for a UI Action in ServiceNow.
type UIAction struct {
	BaseResult
	Name               string `json:"name"`
	Table              string `json:"table"`
	ActionName         string `json:"action_name"`
	Comments           string `json:"comments"`
	Active             bool   `json:"active,string"`
	Script             string `json:"script"`
	Condition          string `json:"condition"`
	FormButton         bool   `json:"form_button,string"`
	FormButtonV2       bool   `json:"form_button_v2,string"`
	FormContextMenu    bool   `json:"form_context_menu,string"`
	FormLink           bool   `json:"form_link,string"`
	FormMenuButtonV2   bool   `json:"form_menu_button_v2,string"`
	ListAction         bool   `json:"list_action,string"`
	ListBannerButton   bool   `json:"list_banner_button,string"`
	ListButton         bool   `json:"list_button,string"`
	ListChoice         bool   `json:"list_choice,string"`
	ListContextMenu    bool   `json:"list_context_menu,string"`
	ListLink           bool   `json:"list_link,string"`
	Client             bool   `json:"client,string"`
	Onclick            string `json:"onclick"`
	Hint               string `json:"hint"`
	Order              int    `json:"order,string"`
	ShowInsert         bool   `json:"show_insert,string"`
	ShowUpdate         bool   `json:"show_update,string"`
	ShowQuery          bool   `json:"show_query,string"`
	ShowMultipleUpdate bool   `json:"show_multiple_update,string"`
}
