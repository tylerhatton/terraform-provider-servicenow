package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const cssIncludeRelationDependencyID = "dependency_id"
const cssIncludeRelationCSSIncludeID = "css_include_id"
const cssIncludeRelationOrder = "order"

// ResourceCSSIncludeRelation is holding the info about the relation between a CSS Include and a widget dependency.
func ResourceCSSIncludeRelation() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_css_include_relation` manages a relation between a CSS include and a widget dependency within ServiceNow.",

		CreateContext: createResourceCSSIncludeRelation,
		ReadContext:   readResourceCSSIncludeRelation,
		UpdateContext: updateResourceCSSIncludeRelation,
		DeleteContext: deleteResourceCSSIncludeRelation,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			cssIncludeRelationDependencyID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sys ID of the widget dependency this CSS include is associated with.",
			},
			cssIncludeRelationCSSIncludeID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sys ID of the CSS include to associate with the widget dependency.",
			},
			cssIncludeRelationOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The load order for the CSS include within the widget dependency.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceCSSIncludeRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	cssIncludeRelation := &client.CSSIncludeRelation{}
	if err := snowClient.GetObject(ctx, client.EndpointCSSIncludeRelation, data.Id(), cssIncludeRelation); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromCSSIncludeRelation(data, cssIncludeRelation)

	return nil
}

func createResourceCSSIncludeRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	cssIncludeRelation := resourceToCSSIncludeRelation(data)
	if err := snowClient.CreateObject(ctx, client.EndpointCSSIncludeRelation, cssIncludeRelation); err != nil {
		return diag.FromErr(err)
	}

	resourceFromCSSIncludeRelation(data, cssIncludeRelation)

	return readResourceCSSIncludeRelation(ctx, data, serviceNowClient)
}

func updateResourceCSSIncludeRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointCSSIncludeRelation, resourceToCSSIncludeRelation(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceCSSIncludeRelation(ctx, data, serviceNowClient)
}

func deleteResourceCSSIncludeRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointCSSIncludeRelation, data.Id()))
}

func resourceFromCSSIncludeRelation(data *schema.ResourceData, cssIncludeRelation *client.CSSIncludeRelation) {
	data.SetId(cssIncludeRelation.ID)
	data.Set(cssIncludeRelationDependencyID, cssIncludeRelation.DependencyID)
	data.Set(cssIncludeRelationCSSIncludeID, cssIncludeRelation.CSSIncludeID)
	data.Set(cssIncludeRelationOrder, cssIncludeRelation.Order)
	data.Set(commonScope, cssIncludeRelation.Scope)
}

func resourceToCSSIncludeRelation(data *schema.ResourceData) *client.CSSIncludeRelation {
	cssIncludeRelation := client.CSSIncludeRelation{
		DependencyID: data.Get(cssIncludeRelationDependencyID).(string),
		CSSIncludeID: data.Get(cssIncludeRelationCSSIncludeID).(string),
		Order:        data.Get(cssIncludeRelationOrder).(int),
	}
	cssIncludeRelation.ID = data.Id()
	cssIncludeRelation.Scope = data.Get(commonScope).(string)
	return &cssIncludeRelation
}
