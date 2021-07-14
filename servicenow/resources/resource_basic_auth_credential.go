package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		Create: createResourceBasicAuthCredential,
		Read:   readResourceBasicAuthCredential,
		Update: updateResourceBasicAuthCredential,
		Delete: deleteResourceBasicAuthCredential,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceBasicAuthCredential(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	basicAuthCredential := &client.BasicAuthCredential{}
	if err := snowClient.GetObject(client.EndpointBasicAuthCredential, data.Id(), basicAuthCredential); err != nil {
		data.SetId("")
		return err
	}

	resourceFromBasicAuthCredential(data, basicAuthCredential)

	return nil
}

func createResourceBasicAuthCredential(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	basicAuthCredential := resourceToBasicAuthCredential(data)
	if err := snowClient.CreateObject(client.EndpointBasicAuthCredential, basicAuthCredential); err != nil {
		return err
	}

	resourceFromBasicAuthCredential(data, basicAuthCredential)

	return readResourceBasicAuthCredential(data, serviceNowClient)
}

func updateResourceBasicAuthCredential(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointBasicAuthCredential, resourceToBasicAuthCredential(data)); err != nil {
		return err
	}

	return readResourceBasicAuthCredential(data, serviceNowClient)
}

func deleteResourceBasicAuthCredential(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return snowClient.DeleteObject(client.EndpointBasicAuthCredential, data.Id())
}

func resourceFromBasicAuthCredential(data *schema.ResourceData, basicAuthCredential *client.BasicAuthCredential) {
	data.SetId(basicAuthCredential.ID)
	data.Set(basicAuthCredentialName, basicAuthCredential.Name)
	data.Set(basicAuthCredentialOrder, basicAuthCredential.Order)
	data.Set(basicAuthCredentialUserName, basicAuthCredential.UserName)
	data.Set(basicAuthCredentialPassword, basicAuthCredential.Password)
	data.Set(basicAuthCredentialCredentialAlias, basicAuthCredential.CredentialAlias)
	data.Set(commonScope, basicAuthCredential.Scope)
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
	basicAuthCredential.Scope = data.Get(commonScope).(string)
	return &basicAuthCredential
}
