package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const basicAuthCredentialName = "name"
const basicAuthCredentialOrder = "order"
const basicAuthCredentialUserName = "username"
const basicAuthCredentialPassword = "password"
const basicAuthCredentialCredentialAlias = "credential_alias"

// ResourceBasicAuthCredential manages a basic authentication credential in ServiceNow.
func ResourceBasicAuthCredential() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_basic_auth_credential` manages a basic auth credential configuration within ServiceNow.",

		CreateContext: createResourceBasicAuthCredential,
		ReadContext:   readResourceBasicAuthCredential,
		UpdateContext: updateResourceBasicAuthCredential,
		DeleteContext: deleteResourceBasicAuthCredential,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			basicAuthCredentialName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of basic auth credential object.",
			},
			basicAuthCredentialOrder: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "100",
				Description: "The order the credential will be use in a credential alias if multiple credentials are applied.",
			},
			basicAuthCredentialUserName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Username assigned to the credential.",
			},
			basicAuthCredentialPassword: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "Password assigned to the credential.",
			},
			basicAuthCredentialCredentialAlias: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sys ID for credential alias the credential is assigned to.",
			},
		},
	}
}

func readResourceBasicAuthCredential(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	basicAuthCredential := &client.BasicAuthCredential{}
	if err := snowClient.GetObject(client.EndpointBasicAuthCredential, data.Id(), basicAuthCredential); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromBasicAuthCredential(data, basicAuthCredential)

	return nil
}

func createResourceBasicAuthCredential(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	basicAuthCredential := resourceToBasicAuthCredential(data)
	if err := snowClient.CreateObject(client.EndpointBasicAuthCredential, basicAuthCredential); err != nil {
		return diag.FromErr(err)
	}

	resourceFromBasicAuthCredential(data, basicAuthCredential)

	return readResourceBasicAuthCredential(ctx, data, serviceNowClient)
}

func updateResourceBasicAuthCredential(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointBasicAuthCredential, resourceToBasicAuthCredential(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceBasicAuthCredential(ctx, data, serviceNowClient)
}

func deleteResourceBasicAuthCredential(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointBasicAuthCredential, data.Id()))
}

func resourceFromBasicAuthCredential(data *schema.ResourceData, basicAuthCredential *client.BasicAuthCredential) {
	data.SetId(basicAuthCredential.ID)
	data.Set(basicAuthCredentialName, basicAuthCredential.Name)
	data.Set(basicAuthCredentialOrder, basicAuthCredential.Order)
	data.Set(basicAuthCredentialUserName, basicAuthCredential.UserName)
	// ServiceNow returns the password in an encrypted form that differs from the
	// plaintext provided by the user. Never overwrite the user-provided password
	// in state with the server-side encrypted value to avoid perpetual drift.
	data.Set(basicAuthCredentialCredentialAlias, basicAuthCredential.CredentialAlias)
}

func resourceToBasicAuthCredential(data *schema.ResourceData) *client.BasicAuthCredential {
	basicAuthCredential := client.BasicAuthCredential{
		Name:            data.Get(basicAuthCredentialName).(string),
		Order:           data.Get(basicAuthCredentialOrder).(string),
		UserName:        data.Get(basicAuthCredentialUserName).(string),
		Password:        data.Get(basicAuthCredentialPassword).(string),
		CredentialAlias: data.Get(basicAuthCredentialCredentialAlias).(string),
	}
	basicAuthCredential.ID = data.Id()
	return &basicAuthCredential
}
