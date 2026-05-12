package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const uiPageName = "name"
const uiPageDescription = "description"
const uiPageCategory = "category"
const uiPageDirect = "direct"
const uiPageClientScript = "client_script"
const uiPageProcessingScript = "processing_script"
const uiPageHTML = "html"
const uiPageEndpoint = "endpoint"

// ResourceUIPage manages a UI Page in ServiceNow.
func ResourceUIPage() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_ui_page` manages a UI Page configuration within ServiceNow.",

		CreateContext: createResourceUIPage,
		ReadContext:   readResourceUIPage,
		UpdateContext: updateResourceUIPage,
		DeleteContext: deleteResourceUIPage,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			uiPageName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the UI Page.",
			},
			uiPageDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the UI Page.",
			},
			uiPageCategory: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "general",
				Description: "Category the UI Page belongs to.",
			},
			uiPageDirect: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this UI Page can be accessed directly via URL without being embedded in the ServiceNow frame.",
			},
			uiPageClientScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client-side JavaScript to be executed on the UI Page.",
			},
			uiPageProcessingScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server-side script to be executed when the UI Page loads.",
			},
			uiPageHTML: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The HTML body of the UI Page.",
			},
			uiPageEndpoint: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL endpoint for accessing this UI Page.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceUIPage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiPage := &client.UIPage{}
	if err := snowClient.GetObject(ctx, client.EndpointUIPage, data.Id(), uiPage); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUIPage(data, uiPage)

	return nil
}

func createResourceUIPage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	uiPage := resourceToUIPage(data)
	if err := snowClient.CreateObject(ctx, client.EndpointUIPage, uiPage); err != nil {
		return diag.FromErr(err)
	}

	resourceFromUIPage(data, uiPage)

	return readResourceUIPage(ctx, data, serviceNowClient)
}

func updateResourceUIPage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointUIPage, resourceToUIPage(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceUIPage(ctx, data, serviceNowClient)
}

func deleteResourceUIPage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointUIPage, data.Id()))
}

func resourceFromUIPage(data *schema.ResourceData, page *client.UIPage) {
	data.SetId(page.ID)
	data.Set(uiPageName, page.Name)
	data.Set(uiPageDescription, page.Description)
	data.Set(uiPageDirect, page.Direct)
	data.Set(uiPageHTML, page.HTML)
	data.Set(uiPageProcessingScript, page.ProcessingScript)
	data.Set(uiPageClientScript, page.ClientScript)
	data.Set(uiPageCategory, page.Category)
	data.Set(uiPageEndpoint, page.Endpoint)
	data.Set(commonProtectionPolicy, page.ProtectionPolicy)
	data.Set(commonScope, page.Scope)
}

func resourceToUIPage(data *schema.ResourceData) *client.UIPage {
	uiPage := client.UIPage{
		Name:             data.Get(uiPageName).(string),
		Description:      data.Get(uiPageDescription).(string),
		Direct:           data.Get(uiPageDirect).(bool),
		HTML:             data.Get(uiPageHTML).(string),
		ProcessingScript: data.Get(uiPageProcessingScript).(string),
		ClientScript:     data.Get(uiPageClientScript).(string),
		Category:         data.Get(uiPageCategory).(string),
	}
	uiPage.ID = data.Id()
	uiPage.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	uiPage.Scope = data.Get(commonScope).(string)
	return &uiPage
}
