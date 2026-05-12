package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const dataLookupName = "name"
const dataLookupTable = "table"
const dataLookupLookupTable = "lookup_table"
const dataLookupActive = "active"
const dataLookupRunOnInsert = "run_on_insert"
const dataLookupRunOnUpdate = "run_on_update"
const dataLookupRunOnFormChange = "run_on_form_change"

// ResourceDataLookup manages a data lookup definition in ServiceNow.
func ResourceDataLookup() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_data_lookup` manages a data lookup definition (dl_definition) within ServiceNow.",

		CreateContext: createResourceDataLookup,
		ReadContext:   readResourceDataLookup,
		UpdateContext: updateResourceDataLookup,
		DeleteContext: deleteResourceDataLookup,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			dataLookupName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the data lookup definition.",
			},
			dataLookupTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Table that the lookup applies to (maps to ServiceNow's `source_table` column).",
			},
			dataLookupLookupTable: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source table for the lookup values (maps to ServiceNow's `matcher_table` column).",
			},
			dataLookupActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this data lookup is enabled.",
			},
			dataLookupRunOnInsert: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the lookup runs when a target record is inserted.",
			},
			dataLookupRunOnUpdate: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the lookup runs when a target record is updated.",
			},
			dataLookupRunOnFormChange: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the lookup runs from the client form on field changes.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceDataLookup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	dataLookup := &client.DataLookup{}
	if err := snowClient.GetObject(ctx, client.EndpointDataLookup, data.Id(), dataLookup); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromDataLookup(data, dataLookup)

	return nil
}

func createResourceDataLookup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	dataLookup := resourceToDataLookup(data)
	if err := snowClient.CreateObject(ctx, client.EndpointDataLookup, dataLookup); err != nil {
		return diag.FromErr(err)
	}

	resourceFromDataLookup(data, dataLookup)

	return readResourceDataLookup(ctx, data, serviceNowClient)
}

func updateResourceDataLookup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointDataLookup, resourceToDataLookup(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceDataLookup(ctx, data, serviceNowClient)
}

func deleteResourceDataLookup(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointDataLookup, data.Id()))
}

func resourceFromDataLookup(data *schema.ResourceData, dataLookup *client.DataLookup) {
	data.SetId(dataLookup.ID)
	data.Set(dataLookupName, dataLookup.Name)
	data.Set(dataLookupTable, dataLookup.Table)
	data.Set(dataLookupLookupTable, dataLookup.LookupTable)
	data.Set(dataLookupActive, dataLookup.Active)
	data.Set(dataLookupRunOnInsert, dataLookup.RunOnInsert)
	data.Set(dataLookupRunOnUpdate, dataLookup.RunOnUpdate)
	data.Set(dataLookupRunOnFormChange, dataLookup.RunOnFormChange)
	data.Set(commonProtectionPolicy, dataLookup.ProtectionPolicy)
	data.Set(commonScope, dataLookup.Scope)
}

func resourceToDataLookup(data *schema.ResourceData) *client.DataLookup {
	dataLookup := client.DataLookup{
		Name:            data.Get(dataLookupName).(string),
		Table:           data.Get(dataLookupTable).(string),
		LookupTable:     data.Get(dataLookupLookupTable).(string),
		Active:          data.Get(dataLookupActive).(bool),
		RunOnInsert:     data.Get(dataLookupRunOnInsert).(bool),
		RunOnUpdate:     data.Get(dataLookupRunOnUpdate).(bool),
		RunOnFormChange: data.Get(dataLookupRunOnFormChange).(bool),
	}
	dataLookup.ID = data.Id()
	dataLookup.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	dataLookup.Scope = data.Get(commonScope).(string)
	return &dataLookup
}
