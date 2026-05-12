package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const widgetID = "identifier"
const widgetName = "name"
const widgetTemplate = "template"
const widgetCSS = "css"
const widgetPublic = "public"
const widgetRoles = "roles"
const widgetLink = "link_function"
const widgetDescription = "description"
const widgetClientScript = "client_script"
const widgetServerScript = "server_script"
const widgetDemoData = "demo_data"
const widgetOptionSchema = "option_schema"
const widgetHasPreview = "has_preview"
const widgetDataTable = "data_table"
const widgetControllerAs = "controller_as"

// ResourceWidget manages a Widget in ServiceNow.
func ResourceWidget() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_widget` manages a Widget configuration within ServiceNow.",

		CreateContext: createResourceWidget,
		ReadContext:   readResourceWidget,
		UpdateContext: updateResourceWidget,
		DeleteContext: deleteResourceWidget,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			widgetID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for the widget used in scripts and portal pages.",
			},
			widgetName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the widget.",
			},
			widgetTemplate: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The HTML template body of the widget.",
			},
			widgetCSS: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "CSS styles applied to the widget.",
			},
			widgetPublic: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this widget is accessible without authentication.",
			},
			widgetRoles: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-separated list of roles required to use the widget.",
			},
			widgetLink: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Link function script executed when the widget is linked.",
			},
			widgetDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the widget.",
			},
			widgetClientScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Client-side AngularJS controller script for the widget.",
			},
			widgetServerScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Server-side script executed when the widget loads.",
			},
			widgetDemoData: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Demo data in JSON format used when previewing the widget.",
			},
			widgetOptionSchema: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "JSON schema defining the configurable options for the widget.",
			},
			widgetHasPreview: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', the widget supports preview mode.",
			},
			widgetDataTable: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The ServiceNow table the widget fetches data from.",
			},
			widgetControllerAs: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "c",
				Description: "The AngularJS controller alias used in the widget template.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceWidget(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	widget := &client.Widget{}
	if err := snowClient.GetObject(client.EndpointWidget, data.Id(), widget); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromWidget(data, widget)

	return nil
}

func createResourceWidget(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	widget := resourceToWidget(data)
	if err := snowClient.CreateObject(client.EndpointWidget, widget); err != nil {
		return diag.FromErr(err)
	}

	resourceFromWidget(data, widget)

	return readResourceWidget(ctx, data, serviceNowClient)
}

func updateResourceWidget(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointWidget, resourceToWidget(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceWidget(ctx, data, serviceNowClient)
}

func deleteResourceWidget(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointWidget, data.Id()))
}

func resourceFromWidget(data *schema.ResourceData, widget *client.Widget) {
	data.SetId(widget.ID)
	data.Set(widgetID, widget.CustomID)
	data.Set(widgetName, widget.Name)
	data.Set(widgetTemplate, widget.Template)
	data.Set(widgetCSS, widget.CSS)
	data.Set(widgetPublic, widget.Public)
	data.Set(widgetRoles, widget.Roles)
	data.Set(widgetLink, widget.Link)
	data.Set(widgetDescription, widget.Description)
	data.Set(widgetClientScript, widget.ClientScript)
	data.Set(widgetServerScript, widget.ServerScript)
	data.Set(widgetDemoData, widget.DemoData)
	data.Set(widgetOptionSchema, widget.OptionSchema)
	data.Set(widgetHasPreview, widget.HasPreview)
	data.Set(widgetDataTable, widget.DataTable)
	data.Set(widgetControllerAs, widget.ControllerAs)
	data.Set(commonProtectionPolicy, widget.ProtectionPolicy)
	data.Set(commonScope, widget.Scope)
}

func resourceToWidget(data *schema.ResourceData) *client.Widget {
	widget := client.Widget{
		CustomID:     data.Get(widgetID).(string),
		Name:         data.Get(widgetName).(string),
		Template:     data.Get(widgetTemplate).(string),
		CSS:          data.Get(widgetCSS).(string),
		Public:       data.Get(widgetPublic).(bool),
		Roles:        data.Get(widgetRoles).(string),
		Link:         data.Get(widgetLink).(string),
		Description:  data.Get(widgetDescription).(string),
		ClientScript: data.Get(widgetClientScript).(string),
		ServerScript: data.Get(widgetServerScript).(string),
		DemoData:     data.Get(widgetDemoData).(string),
		OptionSchema: data.Get(widgetOptionSchema).(string),
		HasPreview:   data.Get(widgetHasPreview).(bool),
		DataTable:    data.Get(widgetDataTable).(string),
		ControllerAs: data.Get(widgetControllerAs).(string),
	}
	widget.ID = data.Id()
	widget.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	widget.Scope = data.Get(commonScope).(string)
	return &widget
}
