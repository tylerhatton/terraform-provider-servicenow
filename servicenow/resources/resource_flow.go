package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const flowName = "name"
const flowDescription = "description"
const flowActive = "active"
const flowTriggerType = "trigger_type"
const flowApplication = "application"
const flowCategory = "category"
const flowInternalName = "internal_name"
const flowStatus = "status"

// ResourceFlow manages a Flow Designer flow record in ServiceNow.
//
// NOTE: The full Flow Designer schema (steps, variables, triggers, snapshots,
// label cache, etc.) is NOT exposed by this resource. Flow Designer flows are
// composed of many related tables and are normally authored interactively in
// the Flow Designer UI. Users should design flows in the UI and either:
//   - reference an existing flow via the `servicenow_flow` data source, or
//   - manage only the high-level metadata of a flow shell here while authoring
//     the steps/triggers themselves in the UI.
func ResourceFlow() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_flow` manages the high-level metadata of a Flow Designer flow record in ServiceNow. " +
			"The full Flow Designer schema (steps, variables, triggers) is not exposed by Terraform; " +
			"design the body of flows in the Flow Designer UI and reference them here.",

		CreateContext: createResourceFlow,
		ReadContext:   readResourceFlow,
		UpdateContext: updateResourceFlow,
		DeleteContext: deleteResourceFlow,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			flowName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the Flow Designer flow.",
			},
			flowDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the flow.",
			},
			flowActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', the flow is active and can run.",
			},
			flowTriggerType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Trigger type for the flow, e.g. record_created, record_updated, record_deleted, scheduled. Note: triggers are normally stored on sys_hub_trigger_instance records associated with the flow; setting this here is informational only.",
			},
			flowApplication: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "global",
				ForceNew:    true,
				Description: "Sys ID of the sys_app application scope the flow belongs to (maps to sys_scope).",
			},
			flowCategory: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "flow",
				Description: "Category of the flow: flow, subflow or action. Maps to ServiceNow's `type` column on sys_hub_flow.",
			},
			flowInternalName: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Auto-generated internal name derived from the display name.",
			},
			flowStatus: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "draft",
				Description: "Publication status of the flow: draft or published.",
			},
		},
	}
}

func readResourceFlow(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	flow := &client.Flow{}
	if err := snowClient.GetObject(ctx, client.EndpointFlow, data.Id(), flow); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromFlow(data, flow)

	return nil
}

func createResourceFlow(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	flow := resourceToFlow(data)
	if err := snowClient.CreateObject(ctx, client.EndpointFlow, flow); err != nil {
		return diag.FromErr(err)
	}

	resourceFromFlow(data, flow)

	return readResourceFlow(ctx, data, serviceNowClient)
}

func updateResourceFlow(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointFlow, resourceToFlow(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceFlow(ctx, data, serviceNowClient)
}

func deleteResourceFlow(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointFlow, data.Id()))
}

func resourceFromFlow(data *schema.ResourceData, flow *client.Flow) {
	data.SetId(flow.ID)
	data.Set(flowName, flow.Name)
	data.Set(flowDescription, flow.Description)
	data.Set(flowActive, flow.Active)
	data.Set(flowTriggerType, flow.TriggerType)
	data.Set(flowApplication, flow.Scope)
	data.Set(flowCategory, flow.Category)
	data.Set(flowInternalName, flow.InternalName)
	data.Set(flowStatus, flow.Status)
}

func resourceToFlow(data *schema.ResourceData) *client.Flow {
	flow := client.Flow{
		Name:         data.Get(flowName).(string),
		Description:  data.Get(flowDescription).(string),
		Active:       data.Get(flowActive).(bool),
		TriggerType:  data.Get(flowTriggerType).(string),
		Category:     data.Get(flowCategory).(string),
		InternalName: data.Get(flowInternalName).(string),
		Status:       data.Get(flowStatus).(string),
	}
	flow.ID = data.Id()
	flow.Scope = data.Get(flowApplication).(string)
	return &flow
}
