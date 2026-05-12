package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const restMessageHeaderName = "name"
const restMessageHeaderValue = "value"
const restMessageHeaderMessageID = "rest_message_id"

// ResourceRestMessageHeader is holding the info about a header to be applied to a REST method.
func ResourceRestMessageHeader() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_rest_message_header` manages a header to be applied to a REST message within ServiceNow.",

		CreateContext: createResourceRestMessageHeader,
		ReadContext:   readResourceRestMessageHeader,
		UpdateContext: updateResourceRestMessageHeader,
		DeleteContext: deleteResourceRestMessageHeader,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			restMessageHeaderName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the header to add to the HTTP request.",
			},
			restMessageHeaderValue: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the header to add to the HTTP request.",
			},
			restMessageHeaderMessageID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The REST message record ID this header will be applied to.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceRestMessageHeader(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMessageHeader := &client.RestMessageHeader{}
	if err := snowClient.GetObject(client.EndpointRestMessageHeader, data.Id(), restMessageHeader); err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromRestMessageHeader(data, restMessageHeader)

	return nil
}

func createResourceRestMessageHeader(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	restMessageHeader := resourceToRestMessageHeader(data)
	if err := snowClient.CreateObject(client.EndpointRestMessageHeader, restMessageHeader); err != nil {
		return diag.FromErr(err)
	}

	resourceFromRestMessageHeader(data, restMessageHeader)

	return readResourceRestMessageHeader(ctx, data, serviceNowClient)
}

func updateResourceRestMessageHeader(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointRestMessageHeader, resourceToRestMessageHeader(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceRestMessageHeader(ctx, data, serviceNowClient)
}

func deleteResourceRestMessageHeader(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointRestMessageHeader, data.Id()))
}

func resourceFromRestMessageHeader(data *schema.ResourceData, restMessageHeader *client.RestMessageHeader) {
	data.SetId(restMessageHeader.ID)
	data.Set(restMessageHeaderName, restMessageHeader.Name)
	data.Set(restMessageHeaderValue, restMessageHeader.Value)
	data.Set(restMessageHeaderMessageID, restMessageHeader.MessageID)
	data.Set(commonScope, restMessageHeader.Scope)
}

func resourceToRestMessageHeader(data *schema.ResourceData) *client.RestMessageHeader {
	restMessageHeader := client.RestMessageHeader{
		Name:      data.Get(restMessageHeaderName).(string),
		Value:     data.Get(restMessageHeaderValue).(string),
		MessageID: data.Get(restMessageHeaderMessageID).(string),
	}
	restMessageHeader.ID = data.Id()
	restMessageHeader.Scope = data.Get(commonScope).(string)
	return &restMessageHeader
}
