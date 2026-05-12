package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const aliasName = "name"
const aliasParentAlias = "parent_alias"
const aliasType = "type"
const aliasConnectionType = "connection_type"
const aliasMultipleActions = "multiple_actions"
const aliasRetryPolicy = "retry_policy"
const aliasConfigurationTemplate = "configuration_template"

// ResourceAlias manages a connection and credential object in ServiceNow.
func ResourceAlias() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_alias` manages a connection and credential object within ServiceNow to provide connection details to a Flow Designer action.",

		CreateContext: createResourceAlias,
		ReadContext:   readResourceAlias,
		UpdateContext: updateResourceAlias,
		DeleteContext: deleteResourceAlias,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			aliasName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of alias object.",
			},
			aliasParentAlias: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys id of parent alias",
			},
			aliasType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "credential",
				Description: "Type of alias. credential or connection",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					warns, errs = validateStringValue(val.(string), key, []string{
						"connection",
						"credential",
					})
					return
				},
			},
			aliasConnectionType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http_connection",
				Description: "Type of connection. http_connection, jdbc_connection, basic_connection, or jms_connection",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					warns, errs = validateStringValue(val.(string), key, []string{
						"http_connection",
						"jdbc_connection",
						"basic_connection",
						"jms_connection",
					})
					return
				},
			},
			aliasMultipleActions: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, connection alias allows multiple active connections",
			},
			aliasRetryPolicy: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys id of default retry policy of connection alias",
			},
			aliasConfigurationTemplate: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys id of configuration template of connection alias",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceAlias(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	alias := &client.Alias{}
	if err := snowClient.GetObject(client.EndpointAlias, data.Id(), alias); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromAlias(data, alias)

	return nil
}

func createResourceAlias(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	alias := resourceToAlias(data)
	if err := snowClient.CreateObject(client.EndpointAlias, alias); err != nil {
		return diag.FromErr(err)
	}

	resourceFromAlias(data, alias)

	return readResourceAlias(ctx, data, serviceNowClient)
}

func updateResourceAlias(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointAlias, resourceToAlias(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceAlias(ctx, data, serviceNowClient)
}

func deleteResourceAlias(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointAlias, data.Id()))
}

func resourceFromAlias(data *schema.ResourceData, alias *client.Alias) {
	data.SetId(alias.ID)
	data.Set(aliasName, alias.Name)
	data.Set(aliasParentAlias, alias.ParentAlias)
	data.Set(aliasType, alias.Type)
	data.Set(aliasConnectionType, alias.ConnectionType)
	data.Set(aliasMultipleActions, alias.MultipleActions)
	data.Set(aliasRetryPolicy, alias.RetryPolicy)
	data.Set(aliasConfigurationTemplate, alias.ConfigurationTemplate)
	data.Set(commonScope, alias.Scope)
}

func resourceToAlias(data *schema.ResourceData) *client.Alias {
	alias := client.Alias{
		Name:                  data.Get(aliasName).(string),
		ParentAlias:           data.Get(aliasParentAlias).(string),
		Type:                  data.Get(aliasType).(string),
		ConnectionType:        data.Get(aliasConnectionType).(string),
		MultipleActions:       data.Get(aliasMultipleActions).(bool),
		RetryPolicy:           data.Get(aliasRetryPolicy).(string),
		ConfigurationTemplate: data.Get(aliasConfigurationTemplate).(string),
	}
	alias.ID = data.Id()
	alias.Scope = data.Get(commonScope).(string)
	return &alias
}
