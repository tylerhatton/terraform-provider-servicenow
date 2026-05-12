package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const contentCSSName = "name"
const contentCSSType = "type"
const contentCSSUrl = "url"
const contentCSSStyle = "style"

// ResourceContentCSS is holding the info about a style sheet to be included.
func ResourceContentCSS() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_content_css` manages a style sheet(CSS) configuration within ServiceNow.",

		CreateContext: createResourceContentCSS,
		ReadContext:   readResourceContentCSS,
		UpdateContext: updateResourceContentCSS,
		DeleteContext: deleteResourceContentCSS,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			contentCSSName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the content management style sheet.",
			},
			contentCSSType: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "local",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					warns, errs = validateStringValue(val.(string), key, []string{"link", "local"})
					return
				},
				Description: "The type of this content management style sheet. Can be 'link' or 'local'.",
			},
			contentCSSUrl: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Used when 'type' is set to 'link'. Must be an absolute URL to an external style sheet file.",
			},
			contentCSSStyle: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Used when 'type' is set to 'local'. The raw CSS content of this style sheet.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceContentCSS(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	contentCSS := &client.ContentCSS{}
	if err := snowClient.GetObject(client.EndpointContentCSS, data.Id(), contentCSS); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromContentCSS(data, contentCSS)

	return nil
}

func createResourceContentCSS(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	contentCSS := resourceToContentCSS(data)
	if err := snowClient.CreateObject(client.EndpointContentCSS, contentCSS); err != nil {
		return diag.FromErr(err)
	}

	resourceFromContentCSS(data, contentCSS)

	return readResourceContentCSS(ctx, data, serviceNowClient)
}

func updateResourceContentCSS(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointContentCSS, resourceToContentCSS(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceContentCSS(ctx, data, serviceNowClient)
}

func deleteResourceContentCSS(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointContentCSS, data.Id()))
}

func resourceFromContentCSS(data *schema.ResourceData, contentCSS *client.ContentCSS) {
	data.SetId(contentCSS.ID)
	data.Set(contentCSSName, contentCSS.Name)
	data.Set(contentCSSType, contentCSS.Type)
	data.Set(contentCSSUrl, contentCSS.URL)
	data.Set(contentCSSStyle, contentCSS.Style)
	data.Set(commonScope, contentCSS.Scope)
}

func resourceToContentCSS(data *schema.ResourceData) *client.ContentCSS {
	contentCSS := client.ContentCSS{
		Name:  data.Get(contentCSSName).(string),
		Type:  data.Get(contentCSSType).(string),
		URL:   data.Get(contentCSSUrl).(string),
		Style: data.Get(contentCSSStyle).(string),
	}
	contentCSS.ID = data.Id()
	contentCSS.Scope = data.Get(commonScope).(string)
	return &contentCSS
}
