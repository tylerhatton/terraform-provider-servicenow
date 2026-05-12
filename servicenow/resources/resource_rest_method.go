package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const restMethodName = "name"
const restMethodMessageID = "rest_message_id"
const restMethodHTTPMethod = "http_method"
const restMethodRestEndpoint = "rest_endpoint"
const restMethodQualifiedName = "qualified_name"

// ResourceRestMethod is holding the info about a REST method to be included in a REST message.
func ResourceRestMethod() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_rest_method` manages a REST method within ServiceNow.",

		CreateContext: createResourceRestMethod,
		ReadContext:   readResourceRestMethod,
		UpdateContext: updateResourceRestMethod,
		DeleteContext: deleteResourceRestMethod,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			restMethodName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique identifier for this HTTP method record.",
			},
			restMethodMessageID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The REST message record ID this method is based on.",
			},
			restMethodHTTPMethod: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The HTTP method this record implements. Can be 'get', 'post', 'put', 'patch' or 'delete'.",
				ValidateFunc: validation.StringInSlice([]string{"get", "post", "put", "patch", "delete"}, false),
			},
			restMethodRestEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The URL of the REST web service provider this method sends requests to. Can contain variables in the format '${variable}'.",
			},
			restMethodQualifiedName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fully qualified name of this REST method including the parent message name.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceRestMethod(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMethod := &client.RestMethod{}
	if err := snowClient.GetObject(ctx, client.EndpointRestMethod, data.Id(), restMethod); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromRestMethod(data, restMethod)

	return nil
}

func createResourceRestMethod(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMethod := resourceToRestMethod(data)
	if err := snowClient.CreateObject(ctx, client.EndpointRestMethod, restMethod); err != nil {
		return diag.FromErr(err)
	}

	resourceFromRestMethod(data, restMethod)

	return readResourceRestMethod(ctx, data, serviceNowClient)
}

func updateResourceRestMethod(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointRestMethod, resourceToRestMethod(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceRestMethod(ctx, data, serviceNowClient)
}

func deleteResourceRestMethod(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointRestMethod, data.Id()))
}

func resourceFromRestMethod(data *schema.ResourceData, restMethod *client.RestMethod) {
	data.SetId(restMethod.ID)
	data.Set(restMethodName, restMethod.Name)
	data.Set(restMethodMessageID, restMethod.MessageID)
	data.Set(restMethodHTTPMethod, restMethod.HTTPMethod)
	data.Set(restMethodRestEndpoint, restMethod.RestEndpoint)
	data.Set(restMethodQualifiedName, restMethod.QualifiedName)
	data.Set(commonScope, restMethod.Scope)
}

func resourceToRestMethod(data *schema.ResourceData) *client.RestMethod {
	restMethod := client.RestMethod{
		Name:               data.Get(restMethodName).(string),
		MessageID:          data.Get(restMethodMessageID).(string),
		HTTPMethod:         data.Get(restMethodHTTPMethod).(string),
		RestEndpoint:       data.Get(restMethodRestEndpoint).(string),
		AuthenticationType: "inherit_from_parent",
	}
	restMethod.ID = data.Id()
	restMethod.Scope = data.Get(commonScope).(string)
	return &restMethod
}
