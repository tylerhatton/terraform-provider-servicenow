package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const scheduledJobName = "name"
const scheduledJobScript = "script"
const scheduledJobRunType = "run_type"
const scheduledJobRunTime = "run_time"
const scheduledJobRunDayOfWeek = "run_dayofweek"
const scheduledJobRunDayOfMonth = "run_dayofmonth"
const scheduledJobRunPeriod = "run_period"
const scheduledJobRunStart = "run_start"
const scheduledJobActive = "active"
const scheduledJobConditional = "conditional"
const scheduledJobCondition = "condition"

// ResourceScheduledJob manages a scheduled script job in ServiceNow.
func ResourceScheduledJob() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_scheduled_job` manages a scheduled background script job (sysauto_script) within ServiceNow.",

		CreateContext: createResourceScheduledJob,
		ReadContext:   readResourceScheduledJob,
		UpdateContext: updateResourceScheduledJob,
		DeleteContext: deleteResourceScheduledJob,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			scheduledJobName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the scheduled job.",
			},
			scheduledJobScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Javascript code to execute on the schedule.",
			},
			scheduledJobRunType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "daily",
				Description:  "How often the job runs. Allowed: daily, weekly, monthly, once, on_demand, business_calendar, periodically.",
				ValidateFunc: validation.StringInSlice([]string{"daily", "weekly", "monthly", "once", "on_demand", "business_calendar", "periodically"}, false),
			},
			scheduledJobRunTime: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Time of day to run the job (e.g. 1970-01-01 08:00:00).",
			},
			scheduledJobRunDayOfWeek: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Day of week (1-7) to run the job when run_type is weekly.",
			},
			scheduledJobRunDayOfMonth: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Day of month (1-31) to run the job when run_type is monthly. ServiceNow defaults this to 1 when omitted.",
			},
			scheduledJobRunPeriod: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Interval between runs when run_type is periodically (e.g. 1970-01-01 00:00:10).",
			},
			scheduledJobRunStart: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Start date and time for the schedule. ServiceNow defaults this to the time of creation when omitted.",
			},
			scheduledJobActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this scheduled job is enabled.",
			},
			scheduledJobConditional: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, evaluate the condition before running the job.",
			},
			scheduledJobCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Conditional expression that must evaluate true for the job to run.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceScheduledJob(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scheduledJob := &client.ScheduledJob{}
	if err := snowClient.GetObject(ctx, client.EndpointScheduledJob, data.Id(), scheduledJob); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromScheduledJob(data, scheduledJob)

	return nil
}

func createResourceScheduledJob(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scheduledJob := resourceToScheduledJob(data)
	if err := snowClient.CreateObject(ctx, client.EndpointScheduledJob, scheduledJob); err != nil {
		return diag.FromErr(err)
	}

	resourceFromScheduledJob(data, scheduledJob)

	return readResourceScheduledJob(ctx, data, serviceNowClient)
}

func updateResourceScheduledJob(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointScheduledJob, resourceToScheduledJob(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceScheduledJob(ctx, data, serviceNowClient)
}

func deleteResourceScheduledJob(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointScheduledJob, data.Id()))
}

func resourceFromScheduledJob(data *schema.ResourceData, scheduledJob *client.ScheduledJob) {
	data.SetId(scheduledJob.ID)
	data.Set(scheduledJobName, scheduledJob.Name)
	data.Set(scheduledJobScript, scheduledJob.Script)
	data.Set(scheduledJobRunType, scheduledJob.RunType)
	data.Set(scheduledJobRunTime, scheduledJob.RunTime)
	data.Set(scheduledJobRunDayOfWeek, scheduledJob.RunDayOfWeek)
	data.Set(scheduledJobRunDayOfMonth, scheduledJob.RunDayOfMonth)
	data.Set(scheduledJobRunPeriod, scheduledJob.RunPeriod)
	data.Set(scheduledJobRunStart, scheduledJob.RunStart)
	data.Set(scheduledJobActive, scheduledJob.Active)
	data.Set(scheduledJobConditional, scheduledJob.Conditional)
	data.Set(scheduledJobCondition, scheduledJob.Condition)
	data.Set(commonProtectionPolicy, scheduledJob.ProtectionPolicy)
	data.Set(commonScope, scheduledJob.Scope)
}

func resourceToScheduledJob(data *schema.ResourceData) *client.ScheduledJob {
	scheduledJob := client.ScheduledJob{
		Name:          data.Get(scheduledJobName).(string),
		Script:        data.Get(scheduledJobScript).(string),
		RunType:       data.Get(scheduledJobRunType).(string),
		RunTime:       data.Get(scheduledJobRunTime).(string),
		RunDayOfWeek:  data.Get(scheduledJobRunDayOfWeek).(string),
		RunDayOfMonth: data.Get(scheduledJobRunDayOfMonth).(string),
		RunPeriod:     data.Get(scheduledJobRunPeriod).(string),
		RunStart:      data.Get(scheduledJobRunStart).(string),
		Active:        data.Get(scheduledJobActive).(bool),
		Conditional:   data.Get(scheduledJobConditional).(bool),
		Condition:     data.Get(scheduledJobCondition).(string),
	}
	scheduledJob.ID = data.Id()
	scheduledJob.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	scheduledJob.Scope = data.Get(commonScope).(string)
	return &scheduledJob
}
