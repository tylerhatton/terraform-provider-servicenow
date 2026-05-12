package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const systemPropertyRelationCategoryID = "category_id"
const systemPropertyRelationPropertyID = "property_id"
const systemPropertyRelationOrder = "order"

// ResourceSystemPropertyRelation manages a System Property in ServiceNow.
func ResourceSystemPropertyRelation() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_system_property_relation` manages a relation between system property and system property category within ServiceNow.",

		CreateContext: createResourceSystemPropertyRelation,
		ReadContext:   readResourceSystemPropertyRelation,
		UpdateContext: updateResourceSystemPropertyRelation,
		DeleteContext: deleteResourceSystemPropertyRelation,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			systemPropertyRelationCategoryID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "System Property Category ID to link.",
			},
			systemPropertyRelationPropertyID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the System Property to link.",
			},
			systemPropertyRelationOrder: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1",
				Description: "The display order of the system property within the category.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceSystemPropertyRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	systemPropertyRelation := &client.SystemPropertyRelation{}
	if err := snowClient.GetObject(ctx, client.EndpointSystemPropertyRelation, data.Id(), systemPropertyRelation); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromSystemPropertyRelation(data, systemPropertyRelation)

	return nil
}

func createResourceSystemPropertyRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	systemPropertyRelation := resourceToSystemPropertyRelation(data)
	if err := snowClient.CreateObject(ctx, client.EndpointSystemPropertyRelation, systemPropertyRelation); err != nil {
		return diag.FromErr(err)
	}

	resourceFromSystemPropertyRelation(data, systemPropertyRelation)

	return readResourceSystemPropertyRelation(ctx, data, serviceNowClient)
}

func updateResourceSystemPropertyRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointSystemPropertyRelation, resourceToSystemPropertyRelation(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceSystemPropertyRelation(ctx, data, serviceNowClient)
}

func deleteResourceSystemPropertyRelation(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointSystemPropertyRelation, data.Id()))
}

func resourceFromSystemPropertyRelation(data *schema.ResourceData, systemPropertyRelation *client.SystemPropertyRelation) {
	data.SetId(systemPropertyRelation.ID)
	data.Set(systemPropertyRelationCategoryID, systemPropertyRelation.CategoryID)
	data.Set(systemPropertyRelationPropertyID, systemPropertyRelation.PropertyID)
	data.Set(systemPropertyRelationOrder, systemPropertyRelation.Order)
	data.Set(commonScope, systemPropertyRelation.Scope)
}

func resourceToSystemPropertyRelation(data *schema.ResourceData) *client.SystemPropertyRelation {
	systemPropertyRelation := client.SystemPropertyRelation{
		CategoryID: data.Get(systemPropertyRelationCategoryID).(string),
		PropertyID: data.Get(systemPropertyRelationPropertyID).(string),
		Order:      data.Get(systemPropertyRelationOrder).(string),
	}
	systemPropertyRelation.ID = data.Id()
	systemPropertyRelation.Scope = data.Get(commonScope).(string)
	return &systemPropertyRelation
}
