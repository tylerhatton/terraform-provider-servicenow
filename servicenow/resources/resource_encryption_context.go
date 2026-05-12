package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const encryptionContextName = "name"
const encryptionContextType = "type"
const encryptionContextEncryptionKey = "encryption_key"
const encryptionContextDescription = "description"
const encryptionContextActive = "active"
const encryptionContextAlgorithm = "algorithm"

// ResourceEncryptionContext manages an encryption context record in ServiceNow.
//
// NOTE: The sys_encryption_context table is part of the ServiceNow Edge
// Encryption feature and is only available on instances where that plugin is
// installed and activated. On other instances ServiceNow returns
// "Invalid table: sys_encryption_context" and this resource cannot be used.
func ResourceEncryptionContext() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_encryption_context` manages an encryption context record within ServiceNow Edge Encryption. " +
			"Requires the Edge Encryption plugin to be active on the target instance.",

		CreateContext: createResourceEncryptionContext,
		ReadContext:   readResourceEncryptionContext,
		UpdateContext: updateResourceEncryptionContext,
		DeleteContext: deleteResourceEncryptionContext,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			encryptionContextName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the encryption context record.",
			},
			encryptionContextType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "standard",
				Description:  "Type of encryption context. Allowed values are standard, equality_preserving, order_preserving or blind.",
				ValidateFunc: validation.StringInSlice([]string{"standard", "equality_preserving", "order_preserving", "blind"}, false),
			},
			encryptionContextEncryptionKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "Encryption key associated with the context. Sensitive: the value is not read back from ServiceNow.",
			},
			encryptionContextDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the encryption context.",
			},
			encryptionContextActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', the encryption context is active.",
			},
			encryptionContextAlgorithm: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "AES_256",
				Description:  "Encryption algorithm to use. Allowed values are AES_128 or AES_256.",
				ValidateFunc: validation.StringInSlice([]string{"AES_128", "AES_256"}, false),
			},
		},
	}
}

func readResourceEncryptionContext(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	encryptionContext := &client.EncryptionContext{}
	if err := snowClient.GetObject(ctx, client.EndpointEncryptionContext, data.Id(), encryptionContext); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromEncryptionContext(data, encryptionContext)

	return nil
}

func createResourceEncryptionContext(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	encryptionContext := resourceToEncryptionContext(data)
	if err := snowClient.CreateObject(ctx, client.EndpointEncryptionContext, encryptionContext); err != nil {
		return diag.FromErr(err)
	}

	resourceFromEncryptionContext(data, encryptionContext)

	return readResourceEncryptionContext(ctx, data, serviceNowClient)
}

func updateResourceEncryptionContext(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointEncryptionContext, resourceToEncryptionContext(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceEncryptionContext(ctx, data, serviceNowClient)
}

func deleteResourceEncryptionContext(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointEncryptionContext, data.Id()))
}

func resourceFromEncryptionContext(data *schema.ResourceData, encryptionContext *client.EncryptionContext) {
	data.SetId(encryptionContext.ID)
	data.Set(encryptionContextName, encryptionContext.Name)
	data.Set(encryptionContextType, encryptionContext.Type)
	// Sensitive encryption_key is not read back to prevent drift against the
	// server-side stored/encrypted value.
	data.Set(encryptionContextDescription, encryptionContext.Description)
	data.Set(encryptionContextActive, encryptionContext.Active)
	data.Set(encryptionContextAlgorithm, encryptionContext.Algorithm)
}

func resourceToEncryptionContext(data *schema.ResourceData) *client.EncryptionContext {
	encryptionContext := client.EncryptionContext{
		Name:          data.Get(encryptionContextName).(string),
		Type:          data.Get(encryptionContextType).(string),
		EncryptionKey: data.Get(encryptionContextEncryptionKey).(string),
		Description:   data.Get(encryptionContextDescription).(string),
		Active:        data.Get(encryptionContextActive).(bool),
		Algorithm:     data.Get(encryptionContextAlgorithm).(string),
	}
	encryptionContext.ID = data.Id()
	return &encryptionContext
}
