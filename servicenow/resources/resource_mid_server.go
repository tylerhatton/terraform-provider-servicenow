package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const midServerName = "name"
const midServerHostName = "host_name"
const midServerStatus = "status"
const midServerVersion = "version"
const midServerValidated = "validated"
const midServerOSName = "os_name"
const midServerOSVersion = "os_version"
const midServerDescription = "description"
const midServerMidUser = "mid_user"
const midServerStarted = "started"
const midServerLastRefresh = "last_refresh"
const midServerAgentType = "agent_type"
const midServerLinuxUserName = "linux_user_name"

// ResourceMidServer manages a MID server (ECC agent) record in ServiceNow.
//
// IMPORTANT: MID servers are typically created by installing the MID server agent
// on a host, which then auto-registers with ServiceNow. Creating a MID server
// record purely through this resource will produce a placeholder row that has no
// running agent behind it. In most workflows you should use the
// `servicenow_mid_server` data source to read an existing MID server registered
// by the installer.
func ResourceMidServer() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_mid_server` manages a MID server (ECC agent) record in ServiceNow. " +
			"MID servers are normally registered by the MID server installer; create here only " +
			"to manage existing records or placeholders.",

		CreateContext: createResourceMidServer,
		ReadContext:   readResourceMidServer,
		UpdateContext: updateResourceMidServer,
		DeleteContext: deleteResourceMidServer,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			midServerName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the MID server record.",
			},
			midServerHostName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Host name the MID server runs on.",
			},
			midServerStatus: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Down",
				Description: "Operational status of the MID server (typically computed from MID heartbeat).",
			},
			midServerVersion: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Version string reported by the running MID server agent.",
			},
			midServerValidated: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', the MID server has been validated.",
			},
			midServerOSName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Operating system distribution running the MID server agent (maps to ServiceNow's host_os_distribution column).",
			},
			midServerOSVersion: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Operating system version on the MID server host (maps to ServiceNow's host_os_version column).",
			},
			midServerDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the MID server. ServiceNow may or may not surface this column depending on customisations.",
			},
			midServerMidUser: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Logged-in user account the MID service runs under (maps to ServiceNow's user_name column).",
			},
			midServerStarted: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Date/time the MID server agent was last started.",
			},
			midServerLastRefresh: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Date/time of the last refresh from the running MID server agent.",
			},
			midServerAgentType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Type/category of the MID server (maps to ServiceNow's type column).",
			},
			midServerLinuxUserName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Linux user name the MID server runs under. May not be present on every ServiceNow instance.",
			},
		},
	}
}

func readResourceMidServer(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	midServer := &client.MidServer{}
	if err := snowClient.GetObject(ctx, client.EndpointMidServer, data.Id(), midServer); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromMidServer(data, midServer)

	return nil
}

func createResourceMidServer(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	midServer := resourceToMidServer(data)
	if err := snowClient.CreateObject(ctx, client.EndpointMidServer, midServer); err != nil {
		return diag.FromErr(err)
	}

	resourceFromMidServer(data, midServer)

	return readResourceMidServer(ctx, data, serviceNowClient)
}

func updateResourceMidServer(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointMidServer, resourceToMidServer(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceMidServer(ctx, data, serviceNowClient)
}

func deleteResourceMidServer(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointMidServer, data.Id()))
}

func resourceFromMidServer(data *schema.ResourceData, midServer *client.MidServer) {
	data.SetId(midServer.ID)
	data.Set(midServerName, midServer.Name)
	data.Set(midServerHostName, midServer.HostName)
	data.Set(midServerStatus, midServer.Status)
	data.Set(midServerVersion, midServer.Version)
	data.Set(midServerValidated, midServer.Validated)
	data.Set(midServerOSName, midServer.OSName)
	data.Set(midServerOSVersion, midServer.OSVersion)
	data.Set(midServerDescription, midServer.Description)
	data.Set(midServerMidUser, midServer.MidUser)
	data.Set(midServerStarted, midServer.Started)
	data.Set(midServerLastRefresh, midServer.LastRefresh)
	data.Set(midServerAgentType, midServer.AgentType)
	data.Set(midServerLinuxUserName, midServer.LinuxUserName)
}

func resourceToMidServer(data *schema.ResourceData) *client.MidServer {
	midServer := client.MidServer{
		Name:          data.Get(midServerName).(string),
		HostName:      data.Get(midServerHostName).(string),
		Status:        data.Get(midServerStatus).(string),
		Version:       data.Get(midServerVersion).(string),
		Validated:     data.Get(midServerValidated).(bool),
		OSName:        data.Get(midServerOSName).(string),
		OSVersion:     data.Get(midServerOSVersion).(string),
		Description:   data.Get(midServerDescription).(string),
		MidUser:       data.Get(midServerMidUser).(string),
		Started:       data.Get(midServerStarted).(string),
		LastRefresh:   data.Get(midServerLastRefresh).(string),
		AgentType:     data.Get(midServerAgentType).(string),
		LinuxUserName: data.Get(midServerLinuxUserName).(string),
	}
	midServer.ID = data.Id()
	return &midServer
}
