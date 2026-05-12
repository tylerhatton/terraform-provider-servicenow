package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const widgetDepRelationDependencyID = "dependency_id"
const widgetDepRelationWidgetID = "widget_id"

// ResourceWidgetDependencyRelation is holding the relationship between a widget and a widget dependency (many-2-many).
func ResourceWidgetDependencyRelation() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_widget_dependency_relation` manages a relationship between widget and widget dependency within ServiceNow.",

		CreateContext: createResourceWidgetDepRelation,
		ReadContext:   readResourceWidgetDepRelation,
		UpdateContext: updateResourceWidgetDepRelation,
		DeleteContext: deleteResourceWidgetDepRelation,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			widgetDepRelationDependencyID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The sys ID of the widget dependency to link.",
			},
			widgetDepRelationWidgetID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The sys ID of the widget to link the dependency to.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceWidgetDepRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	relation := &client.WidgetDependencyRelation{}
	if err := snowClient.GetObject(client.EndpointWidgetDependencyRelation, data.Id(), relation); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromWidgetDepRelation(data, relation)

	return nil
}

func createResourceWidgetDepRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	relation := resourceToWidgetDepRelation(data)
	if err := snowClient.CreateObject(client.EndpointWidgetDependencyRelation, relation); err != nil {
		return diag.FromErr(err)
	}

	resourceFromWidgetDepRelation(data, relation)

	return readResourceWidgetDepRelation(ctx, data, serviceNowClient)
}

func updateResourceWidgetDepRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointWidgetDependencyRelation, resourceToWidgetDepRelation(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceWidgetDepRelation(ctx, data, serviceNowClient)
}

func deleteResourceWidgetDepRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointWidgetDependencyRelation, data.Id()))
}

func resourceFromWidgetDepRelation(data *schema.ResourceData, relation *client.WidgetDependencyRelation) {
	data.SetId(relation.ID)
	data.Set(widgetDepRelationDependencyID, relation.DependencyID)
	data.Set(widgetDepRelationWidgetID, relation.WidgetID)
	data.Set(commonScope, relation.Scope)
}

func resourceToWidgetDepRelation(data *schema.ResourceData) *client.WidgetDependencyRelation {
	relation := client.WidgetDependencyRelation{
		DependencyID: data.Get(widgetDepRelationDependencyID).(string),
		WidgetID:     data.Get(widgetDepRelationWidgetID).(string),
	}
	relation.ID = data.Id()
	relation.Scope = data.Get(commonScope).(string)
	return &relation
}
