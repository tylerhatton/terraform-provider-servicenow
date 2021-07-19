package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const httpConnectionName = "name"
const httpConnectionActive = "active"
const httpConnectionCredential = "credential"
const httpConnectionConnectionAlias = "connection_alias"
const httpConnectionConnectionUrl = "connection_url"
const httpConnectionUseMidServer = "use_mid_server"
const httpConnectionMidSelection = "mid_selection"
const httpConnectionMidServer = "mid_server"

// ResourceHttpConnection manages an HTTP Connection configuration in ServiceNow.
func ResourceHttpConnection() *schema.Resource {
	return &schema.Resource{
		Create: createResourceHttpConnection,
		Read:   readResourceHttpConnection,
		Update: updateResourceHttpConnection,
		Delete: deleteResourceHttpConnection,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			httpConnectionName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the connection configuration.",
			},
			httpConnectionActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', this property will enable the http connection.",
			},
			httpConnectionCredential: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of associated credential alias configuration.",
			},
			httpConnectionConnectionAlias: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sys ID of associated connection alias configuration.",
			},
			httpConnectionConnectionUrl: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Base URL of HTTP connection configuration.",
			},
			httpConnectionUseMidServer: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the HTTP connection server will use a mid server.",
			},
			httpConnectionMidSelection: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "auto_select",
				Description: "Decides which mid server is used. auto_select or specific_mid_server",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					warns, errs = validateStringValue(val.(string), key, []string{
						"auto_select",
						"specific_mid_server",
					})
					return
				},
			},
			httpConnectionMidServer: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of associated mid server if in use.",
			},
		},
	}
}

func readResourceHttpConnection(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	httpConnection := &client.HttpConnection{}
	if err := snowClient.GetObject(client.EndpointHttpConnection, data.Id(), httpConnection); err != nil {
		data.SetId("")
		return err
	}

	resourceFromHttpConnection(data, httpConnection)

	return nil
}

func createResourceHttpConnection(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	httpConnection := resourceToHttpConnection(data)
	if err := snowClient.CreateObject(client.EndpointHttpConnection, httpConnection); err != nil {
		return err
	}

	resourceFromHttpConnection(data, httpConnection)

	return readResourceHttpConnection(data, serviceNowClient)
}

func updateResourceHttpConnection(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointHttpConnection, resourceToHttpConnection(data)); err != nil {
		return err
	}

	return readResourceHttpConnection(data, serviceNowClient)
}

func deleteResourceHttpConnection(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return snowClient.DeleteObject(client.EndpointHttpConnection, data.Id())
}

func resourceFromHttpConnection(data *schema.ResourceData, httpConnection *client.HttpConnection) {
	data.SetId(httpConnection.ID)
	data.Set(httpConnectionName, httpConnection.Name)
	data.Set(httpConnectionActive, httpConnection.Active)
	data.Set(httpConnectionCredential, httpConnection.Credential)
	data.Set(httpConnectionConnectionAlias, httpConnection.ConnectionAlias)
	data.Set(httpConnectionConnectionUrl, httpConnection.ConnectionUrl)
	data.Set(httpConnectionUseMidServer, httpConnection.UseMidServer)
	data.Set(httpConnectionMidSelection, httpConnection.MidSelection)
	data.Set(httpConnectionMidServer, httpConnection.MidServer)
}

func resourceToHttpConnection(data *schema.ResourceData) *client.HttpConnection {
	httpConnection := client.HttpConnection{
		Name:            data.Get(httpConnectionName).(string),
		Active:          data.Get(httpConnectionActive).(bool),
		Credential:      data.Get(httpConnectionCredential).(string),
		ConnectionAlias: data.Get(httpConnectionConnectionAlias).(string),
		ConnectionUrl:   data.Get(httpConnectionConnectionUrl).(string),
		UseMidServer:    data.Get(httpConnectionUseMidServer).(bool),
		MidSelection:    data.Get(httpConnectionMidSelection).(string),
		MidServer:       data.Get(httpConnectionMidServer).(string),
	}
	httpConnection.ID = data.Id()
	return &httpConnection
}
