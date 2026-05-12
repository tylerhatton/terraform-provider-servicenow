package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const certificateName = "name"
const certificateShortDescription = "short_description"
const certificateFormat = "format"
const certificateType = "type"
const certificatePEMCertificate = "pem_certificate"
const certificateKeyStorePassword = "key_store_password"
const certificateKeyStore = "key_store"
const certificateExpiration = "expiration"
const certificateSubject = "subject"
const certificateIssuer = "issuer"
const certificateActive = "active"
const certificateValidFrom = "valid_from"

// ResourceCertificate manages a certificate in ServiceNow.
func ResourceCertificate() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_certificate` manages a X.509 certificate or keystore record within ServiceNow.",

		CreateContext: createResourceCertificate,
		ReadContext:   readResourceCertificate,
		UpdateContext: updateResourceCertificate,
		DeleteContext: deleteResourceCertificate,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			certificateName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the certificate record.",
			},
			certificateShortDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Short description of the certificate record.",
			},
			certificateFormat: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "pem",
				Description:  "Format of the certificate. Allowed values are pem, der, jks or pfx.",
				ValidateFunc: validation.StringInSlice([]string{"pem", "der", "jks", "pfx"}, false),
			},
			certificateType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "trust_store",
				Description:  "Type of the certificate. Allowed values are trust_store, java_key_store or key_pair.",
				ValidateFunc: validation.StringInSlice([]string{"trust_store", "java_key_store", "key_pair"}, false),
			},
			certificatePEMCertificate: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "PEM encoded certificate body. Sensitive: the value is not read back from ServiceNow.",
			},
			certificateKeyStorePassword: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "Password used to unlock the keystore. Sensitive: the value is not read back from ServiceNow.",
			},
			certificateKeyStore: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "Base64-encoded keystore content (used when format is jks or pfx). Sensitive: the value is not read back from ServiceNow.",
			},
			certificateExpiration: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Expiration date of the certificate. Typically computed by ServiceNow from the certificate body.",
			},
			certificateSubject: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Subject DN parsed from the certificate.",
			},
			certificateIssuer: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Issuer DN parsed from the certificate.",
			},
			certificateActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', the certificate is active.",
			},
			certificateValidFrom: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Date from which the certificate is valid. Typically computed by ServiceNow from the certificate body.",
			},
		},
	}
}

func readResourceCertificate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	certificate := &client.Certificate{}
	if err := snowClient.GetObject(ctx, client.EndpointCertificate, data.Id(), certificate); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromCertificate(data, certificate)

	return nil
}

func createResourceCertificate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	certificate := resourceToCertificate(data)
	if err := snowClient.CreateObject(ctx, client.EndpointCertificate, certificate); err != nil {
		return diag.FromErr(err)
	}

	resourceFromCertificate(data, certificate)

	return readResourceCertificate(ctx, data, serviceNowClient)
}

func updateResourceCertificate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointCertificate, resourceToCertificate(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceCertificate(ctx, data, serviceNowClient)
}

func deleteResourceCertificate(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointCertificate, data.Id()))
}

func resourceFromCertificate(data *schema.ResourceData, certificate *client.Certificate) {
	data.SetId(certificate.ID)
	data.Set(certificateName, certificate.Name)
	data.Set(certificateShortDescription, certificate.ShortDescription)
	data.Set(certificateFormat, certificate.Format)
	data.Set(certificateType, certificate.Type)
	// Sensitive fields (pem_certificate, key_store, key_store_password) are not read back
	// to prevent perpetual diffs against server-side stored/encrypted values.
	data.Set(certificateExpiration, certificate.Expiration)
	data.Set(certificateSubject, certificate.Subject)
	data.Set(certificateIssuer, certificate.Issuer)
	data.Set(certificateActive, certificate.Active)
	data.Set(certificateValidFrom, certificate.ValidFrom)
}

func resourceToCertificate(data *schema.ResourceData) *client.Certificate {
	certificate := client.Certificate{
		Name:             data.Get(certificateName).(string),
		ShortDescription: data.Get(certificateShortDescription).(string),
		Format:           data.Get(certificateFormat).(string),
		Type:             data.Get(certificateType).(string),
		PEMCertificate:   data.Get(certificatePEMCertificate).(string),
		KeyStorePassword: data.Get(certificateKeyStorePassword).(string),
		KeyStore:         data.Get(certificateKeyStore).(string),
		Expiration:       data.Get(certificateExpiration).(string),
		Subject:          data.Get(certificateSubject).(string),
		Issuer:           data.Get(certificateIssuer).(string),
		Active:           data.Get(certificateActive).(bool),
		ValidFrom:        data.Get(certificateValidFrom).(string),
	}
	certificate.ID = data.Id()
	return &certificate
}
