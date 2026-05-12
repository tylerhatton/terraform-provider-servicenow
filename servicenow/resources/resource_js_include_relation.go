package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const jsIncludeRelationDependencyID = "dependency_id"
const jsIncludeRelationJsIncludeID = "js_include_id"
const jsIncludeRelationOrder = "order"

// ResourceJsIncludeRelation is holding the info about the relation between a js include and a widget dependency.
func ResourceJsIncludeRelation() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_js_include_relation` manages a relation between a js include and a widget dependency within ServiceNow.",

		CreateContext: createResourceJsIncludeRelation,
		ReadContext:   readResourceJsIncludeRelation,
		UpdateContext: updateResourceJsIncludeRelation,
		DeleteContext: deleteResourceJsIncludeRelation,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			jsIncludeRelationDependencyID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sys ID of the widget dependency this JS include is associated with.",
			},
			jsIncludeRelationJsIncludeID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sys ID of the JS include to associate with the widget dependency.",
			},
			jsIncludeRelationOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The load order for the JS include within the widget dependency.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceJsIncludeRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	jsIncludeRelation := &client.JsIncludeRelation{}
	if err := snowClient.GetObject(client.EndpointJsIncludeRelation, data.Id(), jsIncludeRelation); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromJsIncludeRelation(data, jsIncludeRelation)

	return nil
}

func createResourceJsIncludeRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	jsIncludeRelation := resourceToJsIncludeRelation(data)
	if err := snowClient.CreateObject(client.EndpointJsIncludeRelation, jsIncludeRelation); err != nil {
		return diag.FromErr(err)
	}

	resourceFromJsIncludeRelation(data, jsIncludeRelation)

	return readResourceJsIncludeRelation(ctx, data, serviceNowClient)
}

func updateResourceJsIncludeRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointJsIncludeRelation, resourceToJsIncludeRelation(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceJsIncludeRelation(ctx, data, serviceNowClient)
}

func deleteResourceJsIncludeRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointJsIncludeRelation, data.Id()))
}

func resourceFromJsIncludeRelation(data *schema.ResourceData, jsIncludeRelation *client.JsIncludeRelation) {
	data.SetId(jsIncludeRelation.ID)
	data.Set(jsIncludeRelationDependencyID, jsIncludeRelation.DependencyID)
	data.Set(jsIncludeRelationJsIncludeID, jsIncludeRelation.JsIncludeID)
	data.Set(jsIncludeRelationOrder, jsIncludeRelation.Order)
	data.Set(commonScope, jsIncludeRelation.Scope)
}

func resourceToJsIncludeRelation(data *schema.ResourceData) *client.JsIncludeRelation {
	jsIncludeRelation := client.JsIncludeRelation{
		DependencyID: data.Get(jsIncludeRelationDependencyID).(string),
		JsIncludeID:  data.Get(jsIncludeRelationJsIncludeID).(string),
		Order:        data.Get(jsIncludeRelationOrder).(int),
	}
	jsIncludeRelation.ID = data.Id()
	jsIncludeRelation.Scope = data.Get(commonScope).(string)
	return &jsIncludeRelation
}
