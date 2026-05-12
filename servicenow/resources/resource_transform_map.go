package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const transformMapName = "name"
const transformMapSourceTable = "source_table"
const transformMapTargetTable = "target_table"
const transformMapActive = "active"
const transformMapRunBusinessRules = "run_business_rules"
const transformMapEnforceMandatoryFields = "enforce_mandatory_fields"
const transformMapCopyEmptyFields = "copy_empty_fields"
const transformMapOrder = "order"
const transformMapScript = "script"

// ResourceTransformMap manages a transform map in ServiceNow.
func ResourceTransformMap() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_transform_map` manages a transform map (sys_transform_map) within ServiceNow for import data set field mapping.",

		CreateContext: createResourceTransformMap,
		ReadContext:   readResourceTransformMap,
		UpdateContext: updateResourceTransformMap,
		DeleteContext: deleteResourceTransformMap,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			transformMapName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the transform map.",
			},
			transformMapSourceTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source table (typically an import set staging table).",
			},
			transformMapTargetTable: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target table that the transform map will write records to.",
			},
			transformMapActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this transform map is enabled.",
			},
			transformMapRunBusinessRules: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, run business rules on the target table during transformation.",
			},
			transformMapEnforceMandatoryFields: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "No",
				Description: "Mandatory field enforcement setting (No, error, warn).",
			},
			transformMapCopyEmptyFields: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, copy empty source values to the target record.",
			},
			transformMapOrder: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Order in which the map is executed.",
			},
			transformMapScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Optional transform script to execute during transformation. ServiceNow seeds a default function body when omitted.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceTransformMap(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	transformMap := &client.TransformMap{}
	if err := snowClient.GetObject(ctx, client.EndpointTransformMap, data.Id(), transformMap); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromTransformMap(data, transformMap)

	return nil
}

func createResourceTransformMap(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	transformMap := resourceToTransformMap(data)
	if err := snowClient.CreateObject(ctx, client.EndpointTransformMap, transformMap); err != nil {
		return diag.FromErr(err)
	}

	resourceFromTransformMap(data, transformMap)

	return readResourceTransformMap(ctx, data, serviceNowClient)
}

func updateResourceTransformMap(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointTransformMap, resourceToTransformMap(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceTransformMap(ctx, data, serviceNowClient)
}

func deleteResourceTransformMap(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointTransformMap, data.Id()))
}

func resourceFromTransformMap(data *schema.ResourceData, transformMap *client.TransformMap) {
	data.SetId(transformMap.ID)
	data.Set(transformMapName, transformMap.Name)
	data.Set(transformMapSourceTable, transformMap.SourceTable)
	data.Set(transformMapTargetTable, transformMap.TargetTable)
	data.Set(transformMapActive, transformMap.Active)
	data.Set(transformMapRunBusinessRules, transformMap.RunBusinessRules)
	data.Set(transformMapEnforceMandatoryFields, transformMap.EnforceMandatoryFields)
	data.Set(transformMapCopyEmptyFields, transformMap.CopyEmptyFields)
	data.Set(transformMapOrder, transformMap.Order)
	data.Set(transformMapScript, transformMap.Script)
	data.Set(commonProtectionPolicy, transformMap.ProtectionPolicy)
	data.Set(commonScope, transformMap.Scope)
}

func resourceToTransformMap(data *schema.ResourceData) *client.TransformMap {
	transformMap := client.TransformMap{
		Name:                   data.Get(transformMapName).(string),
		SourceTable:            data.Get(transformMapSourceTable).(string),
		TargetTable:            data.Get(transformMapTargetTable).(string),
		Active:                 data.Get(transformMapActive).(bool),
		RunBusinessRules:       data.Get(transformMapRunBusinessRules).(bool),
		EnforceMandatoryFields: data.Get(transformMapEnforceMandatoryFields).(string),
		CopyEmptyFields:        data.Get(transformMapCopyEmptyFields).(bool),
		Order:                  data.Get(transformMapOrder).(int),
		Script:                 data.Get(transformMapScript).(string),
	}
	transformMap.ID = data.Id()
	transformMap.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	transformMap.Scope = data.Get(commonScope).(string)
	return &transformMap
}
