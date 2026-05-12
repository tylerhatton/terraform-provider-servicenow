package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const uiActionName = "name"
const uiActionTable = "table"
const uiActionActionName = "action_name"
const uiActionComments = "comments"
const uiActionActive = "active"
const uiActionScript = "script"
const uiActionCondition = "condition"
const uiActionFormButton = "form_button"
const uiActionFormButtonV2 = "form_button_v2"
const uiActionFormContextMenu = "form_context_menu"
const uiActionFormLink = "form_link"
const uiActionFormMenuButtonV2 = "form_menu_button_v2"
const uiActionListAction = "list_action"
const uiActionListBannerButton = "list_banner_button"
const uiActionListButton = "list_button"
const uiActionListChoice = "list_choice"
const uiActionListContextMenu = "list_context_menu"
const uiActionListLink = "list_link"
const uiActionClient = "client"
const uiActionOnclick = "onclick"
const uiActionHint = "hint"
const uiActionOrder = "order"
const uiActionShowInsert = "show_insert"
const uiActionShowUpdate = "show_update"
const uiActionShowQuery = "show_query"
const uiActionShowMultipleUpdate = "show_multiple_update"

// ResourceUIAction manages a UI Action in ServiceNow.
func ResourceUIAction() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_ui_action` manages a UI action within ServiceNow. UI actions are buttons, links, and menu items that appear on forms and lists, allowing users to perform specific actions.",

		CreateContext: createResourceUIAction,
		ReadContext:   readResourceUIAction,
		UpdateContext: updateResourceUIAction,
		DeleteContext: deleteResourceUIAction,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			uiActionName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the UI action.",
			},
			uiActionTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the table this UI action applies to. Cannot be changed once created.",
			},
			uiActionActionName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The action name used by the system to identify this UI action.",
			},
			uiActionComments: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comments describing this UI action.",
			},
			uiActionActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this UI action is enabled.",
			},
			uiActionScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Script executed when this UI action runs.",
			},
			uiActionCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Condition under which this UI action is displayed.",
			},
			uiActionFormButton: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears as a button on the form.",
			},
			uiActionFormButtonV2: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears as a v2 form button.",
			},
			uiActionFormContextMenu: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears in the form context menu.",
			},
			uiActionFormLink: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears as a related link on the form.",
			},
			uiActionFormMenuButtonV2: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears in the v2 form menu button.",
			},
			uiActionListAction: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears as a list action.",
			},
			uiActionListBannerButton: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears as a list banner button.",
			},
			uiActionListButton: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears as a list button.",
			},
			uiActionListChoice: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears as a list choice action.",
			},
			uiActionListContextMenu: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears in the list context menu.",
			},
			uiActionListLink: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action appears as a related link on the list.",
			},
			uiActionClient: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this UI action runs client-side script when clicked.",
			},
			uiActionOnclick: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Client-side onclick JavaScript handler.",
			},
			uiActionHint: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Tooltip hint text shown on hover.",
			},
			uiActionOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The order in which this UI action appears relative to others.",
			},
			uiActionShowInsert: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the UI action is visible when the record is being inserted.",
			},
			uiActionShowUpdate: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the UI action is visible when the record is being updated.",
			},
			uiActionShowQuery: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the UI action is visible when querying records.",
			},
			uiActionShowMultipleUpdate: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the UI action is visible during multiple update operations.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceUIAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiAction := &client.UIAction{}
	if err := snowClient.GetObject(ctx, client.EndpointUIAction, data.Id(), uiAction); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIAction(data, uiAction)

	return nil
}

func createResourceUIAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiAction := resourceToUIAction(data)
	if err := snowClient.CreateObject(ctx, client.EndpointUIAction, uiAction); err != nil {
		return diag.FromErr(err)
	}

	resourceFromUIAction(data, uiAction)

	return readResourceUIAction(ctx, data, serviceNowClient)
}

func updateResourceUIAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointUIAction, resourceToUIAction(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceUIAction(ctx, data, serviceNowClient)
}

func deleteResourceUIAction(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointUIAction, data.Id()))
}

func resourceFromUIAction(data *schema.ResourceData, uiAction *client.UIAction) {
	data.SetId(uiAction.ID)
	data.Set(uiActionName, uiAction.Name)
	data.Set(uiActionTable, uiAction.Table)
	data.Set(uiActionActionName, uiAction.ActionName)
	data.Set(uiActionComments, uiAction.Comments)
	data.Set(uiActionActive, uiAction.Active)
	data.Set(uiActionScript, uiAction.Script)
	data.Set(uiActionCondition, uiAction.Condition)
	data.Set(uiActionFormButton, uiAction.FormButton)
	data.Set(uiActionFormButtonV2, uiAction.FormButtonV2)
	data.Set(uiActionFormContextMenu, uiAction.FormContextMenu)
	data.Set(uiActionFormLink, uiAction.FormLink)
	data.Set(uiActionFormMenuButtonV2, uiAction.FormMenuButtonV2)
	data.Set(uiActionListAction, uiAction.ListAction)
	data.Set(uiActionListBannerButton, uiAction.ListBannerButton)
	data.Set(uiActionListButton, uiAction.ListButton)
	data.Set(uiActionListChoice, uiAction.ListChoice)
	data.Set(uiActionListContextMenu, uiAction.ListContextMenu)
	data.Set(uiActionListLink, uiAction.ListLink)
	data.Set(uiActionClient, uiAction.Client)
	data.Set(uiActionOnclick, uiAction.Onclick)
	data.Set(uiActionHint, uiAction.Hint)
	data.Set(uiActionOrder, uiAction.Order)
	data.Set(uiActionShowInsert, uiAction.ShowInsert)
	data.Set(uiActionShowUpdate, uiAction.ShowUpdate)
	data.Set(uiActionShowQuery, uiAction.ShowQuery)
	data.Set(uiActionShowMultipleUpdate, uiAction.ShowMultipleUpdate)
	data.Set(commonProtectionPolicy, uiAction.ProtectionPolicy)
	data.Set(commonScope, uiAction.Scope)
}

func resourceToUIAction(data *schema.ResourceData) *client.UIAction {
	uiAction := client.UIAction{
		Name:               data.Get(uiActionName).(string),
		Table:              data.Get(uiActionTable).(string),
		ActionName:         data.Get(uiActionActionName).(string),
		Comments:           data.Get(uiActionComments).(string),
		Active:             data.Get(uiActionActive).(bool),
		Script:             data.Get(uiActionScript).(string),
		Condition:          data.Get(uiActionCondition).(string),
		FormButton:         data.Get(uiActionFormButton).(bool),
		FormButtonV2:       data.Get(uiActionFormButtonV2).(bool),
		FormContextMenu:    data.Get(uiActionFormContextMenu).(bool),
		FormLink:           data.Get(uiActionFormLink).(bool),
		FormMenuButtonV2:   data.Get(uiActionFormMenuButtonV2).(bool),
		ListAction:         data.Get(uiActionListAction).(bool),
		ListBannerButton:   data.Get(uiActionListBannerButton).(bool),
		ListButton:         data.Get(uiActionListButton).(bool),
		ListChoice:         data.Get(uiActionListChoice).(bool),
		ListContextMenu:    data.Get(uiActionListContextMenu).(bool),
		ListLink:           data.Get(uiActionListLink).(bool),
		Client:             data.Get(uiActionClient).(bool),
		Onclick:            data.Get(uiActionOnclick).(string),
		Hint:               data.Get(uiActionHint).(string),
		Order:              data.Get(uiActionOrder).(int),
		ShowInsert:         data.Get(uiActionShowInsert).(bool),
		ShowUpdate:         data.Get(uiActionShowUpdate).(bool),
		ShowQuery:          data.Get(uiActionShowQuery).(bool),
		ShowMultipleUpdate: data.Get(uiActionShowMultipleUpdate).(bool),
	}
	uiAction.ID = data.Id()
	uiAction.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	uiAction.Scope = data.Get(commonScope).(string)
	return &uiAction
}
