package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const restMethodHeaderName = "name"
const restMethodHeaderValue = "value"
const restMethodHeaderMethodID = "rest_method_id"

// ResourceRestMethodHeader is holding the info about a header to be applied to a REST method.
func ResourceRestMethodHeader() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_rest_method_header` manages a header to be applied to a REST method within ServiceNow.",

		CreateContext: createResourceRestMethodHeader,
		ReadContext:   readResourceRestMethodHeader,
		UpdateContext: updateResourceRestMethodHeader,
		DeleteContext: deleteResourceRestMethodHeader,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			restMethodHeaderName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the header to add to the HTTP request.",
			},
			restMethodHeaderValue: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the header to add to the HTTP request.",
			},
			restMethodHeaderMethodID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The REST method record ID this header will be applied to.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceRestMethodHeader(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMethodHeader := &client.RestMethodHeader{}
	if err := snowClient.GetObject(ctx, client.EndpointRestMethodHeader, data.Id(), restMethodHeader); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromRestMethodHeader(data, restMethodHeader)

	return nil
}

func createResourceRestMethodHeader(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMethodHeader := resourceToRestMethodHeader(data)
	if err := snowClient.CreateObject(ctx, client.EndpointRestMethodHeader, restMethodHeader); err != nil {
		return diag.FromErr(err)
	}

	resourceFromRestMethodHeader(data, restMethodHeader)

	return readResourceRestMethodHeader(ctx, data, serviceNowClient)
}

func updateResourceRestMethodHeader(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointRestMethodHeader, resourceToRestMethodHeader(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceRestMethodHeader(ctx, data, serviceNowClient)
}

func deleteResourceRestMethodHeader(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointRestMethodHeader, data.Id()))
}

func resourceFromRestMethodHeader(data *schema.ResourceData, restMethodHeader *client.RestMethodHeader) {
	data.SetId(restMethodHeader.ID)
	data.Set(restMethodHeaderName, restMethodHeader.Name)
	data.Set(restMethodHeaderValue, restMethodHeader.Value)
	data.Set(restMethodHeaderMethodID, restMethodHeader.MethodID)
	data.Set(commonScope, restMethodHeader.Scope)
}

func resourceToRestMethodHeader(data *schema.ResourceData) *client.RestMethodHeader {
	restMethodHeader := client.RestMethodHeader{
		Name:     data.Get(restMethodHeaderName).(string),
		Value:    data.Get(restMethodHeaderValue).(string),
		MethodID: data.Get(restMethodHeaderMethodID).(string),
	}
	restMethodHeader.ID = data.Id()
	restMethodHeader.Scope = data.Get(commonScope).(string)
	return &restMethodHeader
}
