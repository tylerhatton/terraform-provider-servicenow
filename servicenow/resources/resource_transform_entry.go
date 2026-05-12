package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const transformEntryMap = "map"
const transformEntrySourceField = "source_field"
const transformEntryTargetField = "target_field"
const transformEntryCoalesce = "coalesce"
const transformEntryUseSourceScript = "use_source_script"
const transformEntrySourceScript = "source_script"
const transformEntryReferencedValueFieldName = "referenced_value_field_name"

// ResourceTransformEntry manages a transform entry in ServiceNow.
func ResourceTransformEntry() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_transform_entry` manages an individual field mapping (sys_transform_entry) inside a transform map.",

		CreateContext: createResourceTransformEntry,
		ReadContext:   readResourceTransformEntry,
		UpdateContext: updateResourceTransformEntry,
		DeleteContext: deleteResourceTransformEntry,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			transformEntryMap: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "sys_id of the parent transform map this entry belongs to.",
			},
			transformEntrySourceField: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Source field name on the source table.",
			},
			transformEntryTargetField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target field name on the target table.",
			},
			transformEntryCoalesce: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, use this field to match existing records during transformation.",
			},
			transformEntryUseSourceScript: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, use the source script instead of the source field for the value.",
			},
			transformEntrySourceScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Optional source script to compute the value when use_source_script is true. ServiceNow seeds a default function body when the entry is created.",
			},
			transformEntryReferencedValueFieldName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of the field on the reference table to match against (maps to ServiceNow's `reference_value_field` column).",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceTransformEntry(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	transformEntry := &client.TransformEntry{}
	if err := snowClient.GetObject(ctx, client.EndpointTransformEntry, data.Id(), transformEntry); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromTransformEntry(data, transformEntry)

	return nil
}

func createResourceTransformEntry(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	transformEntry := resourceToTransformEntry(data)
	if err := snowClient.CreateObject(ctx, client.EndpointTransformEntry, transformEntry); err != nil {
		return diag.FromErr(err)
	}

	resourceFromTransformEntry(data, transformEntry)

	return readResourceTransformEntry(ctx, data, serviceNowClient)
}

func updateResourceTransformEntry(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointTransformEntry, resourceToTransformEntry(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceTransformEntry(ctx, data, serviceNowClient)
}

func deleteResourceTransformEntry(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointTransformEntry, data.Id()))
}

func resourceFromTransformEntry(data *schema.ResourceData, transformEntry *client.TransformEntry) {
	data.SetId(transformEntry.ID)
	data.Set(transformEntryMap, transformEntry.Map)
	data.Set(transformEntrySourceField, transformEntry.SourceField)
	data.Set(transformEntryTargetField, transformEntry.TargetField)
	data.Set(transformEntryCoalesce, transformEntry.Coalesce)
	data.Set(transformEntryUseSourceScript, transformEntry.UseSourceScript)
	data.Set(transformEntrySourceScript, transformEntry.SourceScript)
	data.Set(transformEntryReferencedValueFieldName, transformEntry.ReferencedValueFieldName)
	data.Set(commonProtectionPolicy, transformEntry.ProtectionPolicy)
	data.Set(commonScope, transformEntry.Scope)
}

func resourceToTransformEntry(data *schema.ResourceData) *client.TransformEntry {
	transformEntry := client.TransformEntry{
		Map:                      data.Get(transformEntryMap).(string),
		SourceField:              data.Get(transformEntrySourceField).(string),
		TargetField:              data.Get(transformEntryTargetField).(string),
		Coalesce:                 data.Get(transformEntryCoalesce).(bool),
		UseSourceScript:          data.Get(transformEntryUseSourceScript).(bool),
		SourceScript:             data.Get(transformEntrySourceScript).(string),
		ReferencedValueFieldName: data.Get(transformEntryReferencedValueFieldName).(string),
	}
	transformEntry.ID = data.Id()
	transformEntry.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	transformEntry.Scope = data.Get(commonScope).(string)
	return &transformEntry
}
