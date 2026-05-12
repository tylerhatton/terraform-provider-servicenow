package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const jdbcConnectionName = "name"
const jdbcConnectionActive = "active"
const jdbcConnectionCredential = "credential"
const jdbcConnectionConnectionAlias = "connection_alias"
const jdbcConnectionConnectionUrl = "connection_url"
const jdbcConnectionDatabaseName = "database_name"
const jdbcConnectionDatabaseType = "database_type"
const jdbcConnectionUseMidServer = "use_mid_server"
const jdbcConnectionMidSelection = "mid_selection"
const jdbcConnectionMidServer = "mid_server"

// ResourceJdbcConnection manages a JDBC Connection configuration in ServiceNow.
func ResourceJdbcConnection() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_jdbc_connection` manages a JDBC connection configuration within ServiceNow. JDBC connections extend the generic connection table and are used to connect ServiceNow to remote relational databases, optionally through a MID server.",

		CreateContext: createResourceJdbcConnection,
		ReadContext:   readResourceJdbcConnection,
		UpdateContext: updateResourceJdbcConnection,
		DeleteContext: deleteResourceJdbcConnection,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			jdbcConnectionName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the JDBC connection configuration.",
			},
			jdbcConnectionActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', this property will enable the JDBC connection.",
			},
			jdbcConnectionCredential: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the associated credential record (typically a JDBC credential alias).",
			},
			jdbcConnectionConnectionAlias: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sys ID of the associated connection alias configuration (type=jdbc_connection).",
			},
			jdbcConnectionConnectionUrl: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JDBC URL for the connection, e.g. jdbc:mysql://host:3306/dbname.",
			},
			jdbcConnectionDatabaseName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Database/schema name to connect to.",
			},
			jdbcConnectionDatabaseType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Database type/driver, e.g. mysql, postgres, oracle, sqlserver. Maps to ServiceNow's jdbc_driver column.",
			},
			jdbcConnectionUseMidServer: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the JDBC connection will be routed through a MID server.",
			},
			jdbcConnectionMidSelection: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "auto_select",
				Description: "Decides which MID server is used. auto_select or specific_mid_server.",
			},
			jdbcConnectionMidServer: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of associated MID server if a specific MID server is selected.",
			},
		},
	}
}

func readResourceJdbcConnection(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	jdbcConnection := &client.JdbcConnection{}
	if err := snowClient.GetObject(ctx, client.EndpointJdbcConnection, data.Id(), jdbcConnection); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromJdbcConnection(data, jdbcConnection)

	return nil
}

func createResourceJdbcConnection(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	jdbcConnection := resourceToJdbcConnection(data)
	if err := snowClient.CreateObject(ctx, client.EndpointJdbcConnection, jdbcConnection); err != nil {
		return diag.FromErr(err)
	}

	resourceFromJdbcConnection(data, jdbcConnection)

	return readResourceJdbcConnection(ctx, data, serviceNowClient)
}

func updateResourceJdbcConnection(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointJdbcConnection, resourceToJdbcConnection(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceJdbcConnection(ctx, data, serviceNowClient)
}

func deleteResourceJdbcConnection(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointJdbcConnection, data.Id()))
}

func resourceFromJdbcConnection(data *schema.ResourceData, jdbcConnection *client.JdbcConnection) {
	data.SetId(jdbcConnection.ID)
	data.Set(jdbcConnectionName, jdbcConnection.Name)
	data.Set(jdbcConnectionActive, jdbcConnection.Active)
	data.Set(jdbcConnectionCredential, jdbcConnection.Credential)
	data.Set(jdbcConnectionConnectionAlias, jdbcConnection.ConnectionAlias)
	data.Set(jdbcConnectionConnectionUrl, jdbcConnection.ConnectionUrl)
	data.Set(jdbcConnectionDatabaseName, jdbcConnection.DatabaseName)
	data.Set(jdbcConnectionDatabaseType, jdbcConnection.DatabaseType)
	data.Set(jdbcConnectionUseMidServer, jdbcConnection.UseMidServer)
	data.Set(jdbcConnectionMidSelection, jdbcConnection.MidSelection)
	data.Set(jdbcConnectionMidServer, jdbcConnection.MidServer)
}

func resourceToJdbcConnection(data *schema.ResourceData) *client.JdbcConnection {
	jdbcConnection := client.JdbcConnection{
		Name:            data.Get(jdbcConnectionName).(string),
		Active:          data.Get(jdbcConnectionActive).(bool),
		Credential:      data.Get(jdbcConnectionCredential).(string),
		ConnectionAlias: data.Get(jdbcConnectionConnectionAlias).(string),
		ConnectionUrl:   data.Get(jdbcConnectionConnectionUrl).(string),
		DatabaseName:    data.Get(jdbcConnectionDatabaseName).(string),
		DatabaseType:    data.Get(jdbcConnectionDatabaseType).(string),
		UseMidServer:    data.Get(jdbcConnectionUseMidServer).(bool),
		MidSelection:    data.Get(jdbcConnectionMidSelection).(string),
		MidServer:       data.Get(jdbcConnectionMidServer).(string),
	}
	jdbcConnection.ID = data.Id()
	return &jdbcConnection
}
