package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const restMessageName = "name"
const restMessageDescription = "description"
const restMessageRestEndpoint = "rest_endpoint"
const restMessageAccess = "access"

// ResourceRestMessage is holding the info about a REST message configuration to be included.
func ResourceRestMessage() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_rest_message` manages a REST message within ServiceNow.",

		CreateContext: createResourceRestMessage,
		ReadContext:   readResourceRestMessage,
		UpdateContext: updateResourceRestMessage,
		DeleteContext: deleteResourceRestMessage,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			restMessageName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Descriptive name for this REST message.",
			},
			restMessageRestEndpoint: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
				Description:  "The URL of the REST web service provider this REST message sends requests to.  Can contain variables in the format '${variable}'.",
			},
			restMessageDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description for this REST message.",
			},
			restMessageAccess: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "package_private",
				Description:  "Whether this REST message can be accessed from only this application scope or all application scopes. Values can be 'package_private' or 'public'.",
				ValidateFunc: validation.StringInSlice([]string{"package_private", "public"}, false),
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceRestMessage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMessage := &client.RestMessage{}
	if err := snowClient.GetObject(ctx, client.EndpointRestMessage, data.Id(), restMessage); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromRestMessage(data, restMessage)

	return nil
}

func createResourceRestMessage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMessage := resourceToRestMessage(data)
	if err := snowClient.CreateObject(ctx, client.EndpointRestMessage, restMessage); err != nil {
		return diag.FromErr(err)
	}

	resourceFromRestMessage(data, restMessage)

	return readResourceRestMessage(ctx, data, serviceNowClient)
}

func updateResourceRestMessage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointRestMessage, resourceToRestMessage(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceRestMessage(ctx, data, serviceNowClient)
}

func deleteResourceRestMessage(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointRestMessage, data.Id()))
}

func resourceFromRestMessage(data *schema.ResourceData, restMessage *client.RestMessage) {
	data.SetId(restMessage.ID)
	data.Set(restMessageName, restMessage.Name)
	data.Set(restMessageDescription, restMessage.Description)
	data.Set(restMessageRestEndpoint, restMessage.RestEndpoint)
	data.Set(restMessageAccess, restMessage.Access)
	data.Set(commonScope, restMessage.Scope)
}

func resourceToRestMessage(data *schema.ResourceData) *client.RestMessage {
	restMessage := client.RestMessage{
		Name:               data.Get(restMessageName).(string),
		Description:        data.Get(restMessageDescription).(string),
		RestEndpoint:       data.Get(restMessageRestEndpoint).(string),
		Access:             data.Get(restMessageAccess).(string),
		AuthenticationType: "no_authentication",
	}
	restMessage.ID = data.Id()
	restMessage.Scope = data.Get(commonScope).(string)
	return &restMessage
}
