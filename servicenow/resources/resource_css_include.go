package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const cssIncludeSource = "source"
const cssIncludeName = "name"
const cssIncludeURL = "url"
const cssIncludeStyleSheetID = "style_sheet_id"

// ResourceCSSInclude is holding the info about a cascading style sheet to be included.
func ResourceCSSInclude() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_css_include` manages a cascading style sheet(CSS) within ServiceNow.",

		CreateContext: createResourceCSSInclude,
		ReadContext:   readResourceCSSInclude,
		UpdateContext: updateResourceCSSInclude,
		DeleteContext: deleteResourceCSSInclude,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			cssIncludeSource: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "url",
				ValidateFunc: validation.StringInSlice([]string{"url", "local"}, false),
				Description:  "Source type of the CSS include. Can be 'url' for an external link or 'local' for a service portal style sheet.",
			},
			cssIncludeName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the CSS include.",
			},
			cssIncludeURL: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.Any(validation.StringIsEmpty, validation.IsURLWithScheme([]string{"http", "https"})),
				Description:  "URL of the external CSS file when source is set to 'url'.",
			},
			cssIncludeStyleSheetID: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The ID of the service portal style sheet to include.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceCSSInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	cssInclude := &client.CSSInclude{}
	if err := snowClient.GetObject(ctx, client.EndpointCSSInclude, data.Id(), cssInclude); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromCSSInclude(data, cssInclude)

	return nil
}

func createResourceCSSInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	cssInclude := resourceToCSSInclude(data)
	if err := snowClient.CreateObject(ctx, client.EndpointCSSInclude, cssInclude); err != nil {
		return diag.FromErr(err)
	}

	resourceFromCSSInclude(data, cssInclude)

	return readResourceCSSInclude(ctx, data, serviceNowClient)
}

func updateResourceCSSInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointCSSInclude, resourceToCSSInclude(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceCSSInclude(ctx, data, serviceNowClient)
}

func deleteResourceCSSInclude(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointCSSInclude, data.Id()))
}

func resourceFromCSSInclude(data *schema.ResourceData, cssInclude *client.CSSInclude) {
	data.SetId(cssInclude.ID)
	data.Set(cssIncludeSource, cssInclude.Source)
	data.Set(cssIncludeName, cssInclude.Name)
	data.Set(cssIncludeURL, cssInclude.URL)
	data.Set(cssIncludeStyleSheetID, cssInclude.StyleSheetID)
	data.Set(commonScope, cssInclude.Scope)
}

func resourceToCSSInclude(data *schema.ResourceData) *client.CSSInclude {
	cssInclude := client.CSSInclude{
		Source:       data.Get(cssIncludeSource).(string),
		Name:         data.Get(cssIncludeName).(string),
		URL:          data.Get(cssIncludeURL).(string),
		StyleSheetID: data.Get(cssIncludeStyleSheetID).(string),
	}
	cssInclude.ID = data.Id()
	cssInclude.Scope = data.Get(commonScope).(string)
	return &cssInclude
}
