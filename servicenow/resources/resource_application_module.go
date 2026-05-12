package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const applicationModuleTitle = "title"
const applicationModuleMenuID = "application_menu_id"
const applicationModuleHint = "hint"
const applicationModuleOrder = "order"
const applicationModuleRoles = "roles"
const applicationModuleActive = "active"
const applicationModuleOverrideRoles = "override_menu_roles"
const applicationModuleLinkType = "link_type"
const applicationModuleLinkArguments = "arguments"
const applicationModuleWindowName = "window_name"
const applicationModuleTableName = "table_name"

// ResourceApplicationModule is a single link in the application navigator.
func ResourceApplicationModule() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_application_module` manages an application menu within ServiceNow creating a link in the application navigator.",

		CreateContext: createResourceApplicationModule,
		ReadContext:   readResourceApplicationModule,
		UpdateContext: updateResourceApplicationModule,
		DeleteContext: deleteResourceApplicationModule,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			applicationModuleTitle: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the module in the application navigator.",
			},
			applicationModuleMenuID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The application Menu ID where this module should reside.",
			},
			applicationModuleHint: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Defines the text that appears in a tooltip when a user points to this module.",
			},
			applicationModuleOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The display order for the module in the application menu.",
			},
			applicationModuleRoles: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-separated list of Roles (names) that can view this application module.",
			},
			applicationModuleActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this application module is in enabled.",
			},
			applicationModuleOverrideRoles: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Show this module when the user has the specified roles. Otherwise the user must have the roles specified by both the application menu and the module.",
			},
			applicationModuleLinkType: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Type of device where this menu will appear. Can be 'DIRECT' for a UI page link or 'LIST' for a table link.",
				ValidateFunc: validation.StringInSlice([]string{"DIRECT", "LIST"}, false),
			},
			applicationModuleLinkArguments: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full name of the UI Page where this module will redirect when link type is 'DIRECT'. When the link type is 'LIST', this is optional.",
			},
			applicationModuleWindowName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the browser window when clicking on a link. For example '_blank' can create a new tab.",
			},
			applicationModuleTableName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The full name of the table where this module will redirect when the link type is 'LIST'.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceApplicationModule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	applicationModule := &client.ApplicationModule{}
	if err := snowClient.GetObject(ctx, client.EndpointApplicationModule, data.Id(), applicationModule); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromApplicationModule(data, applicationModule)

	return nil
}

func createResourceApplicationModule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	applicationModule := resourceToApplicationModule(data)
	if err := snowClient.CreateObject(ctx, client.EndpointApplicationModule, applicationModule); err != nil {
		return diag.FromErr(err)
	}

	resourceFromApplicationModule(data, applicationModule)

	return readResourceApplicationModule(ctx, data, serviceNowClient)
}

func updateResourceApplicationModule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointApplicationModule, resourceToApplicationModule(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceApplicationModule(ctx, data, serviceNowClient)
}

func deleteResourceApplicationModule(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointApplicationModule, data.Id()))
}

func resourceFromApplicationModule(data *schema.ResourceData, applicationModule *client.ApplicationModule) {
	data.SetId(applicationModule.ID)
	data.Set(applicationModuleTitle, applicationModule.Title)
	data.Set(applicationModuleMenuID, applicationModule.MenuID)
	data.Set(applicationModuleHint, applicationModule.Hint)
	data.Set(applicationModuleOrder, applicationModule.Order)
	data.Set(applicationModuleRoles, applicationModule.Roles)
	data.Set(applicationModuleActive, applicationModule.Active)
	data.Set(applicationModuleOverrideRoles, applicationModule.OverrideMenuRoles)
	data.Set(applicationModuleLinkType, applicationModule.LinkType)
	data.Set(applicationModuleLinkArguments, applicationModule.Arguments)
	data.Set(applicationModuleWindowName, applicationModule.WindowName)
	data.Set(applicationModuleTableName, applicationModule.TableName)
	data.Set(commonProtectionPolicy, applicationModule.ProtectionPolicy)
	data.Set(commonScope, applicationModule.Scope)
}

func resourceToApplicationModule(data *schema.ResourceData) *client.ApplicationModule {
	applicationModule := client.ApplicationModule{
		Title:             data.Get(applicationModuleTitle).(string),
		MenuID:            data.Get(applicationModuleMenuID).(string),
		Hint:              data.Get(applicationModuleHint).(string),
		Order:             data.Get(applicationModuleOrder).(int),
		Roles:             data.Get(applicationModuleRoles).(string),
		Active:            data.Get(applicationModuleActive).(bool),
		OverrideMenuRoles: data.Get(applicationModuleOverrideRoles).(bool),
		LinkType:          data.Get(applicationModuleLinkType).(string),
		Arguments:         data.Get(applicationModuleLinkArguments).(string),
		WindowName:        data.Get(applicationModuleWindowName).(string),
		TableName:         data.Get(applicationModuleTableName).(string),
	}
	applicationModule.ID = data.Id()
	applicationModule.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	applicationModule.Scope = data.Get(commonScope).(string)
	return &applicationModule
}
