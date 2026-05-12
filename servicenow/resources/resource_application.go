package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const applicationName = "name"
const applicationScope = "scope"
const applicationVersion = "version"

// ResourceApplication manages an Application in ServiceNow.
func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_application` manages an application within ServiceNow.",

		CreateContext: createResourceApplication,
		ReadContext:   readResourceApplication,
		UpdateContext: updateResourceApplication,
		DeleteContext: deleteResourceApplication,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			applicationName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Application to retrieve from the ServiceNow instance.",
			},
			applicationScope: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique scope of the application. Normally in the format x_[companyCode]_[shortAppId]. Cannot be changed once the application is created.",
			},
			applicationVersion: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1.0.0",
				Description: "The version of the application in semver format.",
			},
		},
	}
}

func readResourceApplication(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	application := &client.Application{}
	if err := snowClient.GetObject(client.EndpointApplication, data.Id(), application); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromApplication(data, application)

	return nil
}

func createResourceApplication(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	application := resourceToApplication(data)
	if err := snowClient.CreateObject(client.EndpointApplication, application); err != nil {
		return diag.FromErr(err)
	}

	resourceFromApplication(data, application)

	return readResourceApplication(ctx, data, serviceNowClient)
}

func updateResourceApplication(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointApplication, resourceToApplication(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceApplication(ctx, data, serviceNowClient)
}

func deleteResourceApplication(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(client.EndpointApplication, data.Id()))
}

func resourceFromApplication(data *schema.ResourceData, application *client.Application) {
	data.SetId(application.ID)
	data.Set(applicationName, application.Name)
	data.Set(applicationScope, application.Scope)
	data.Set(applicationVersion, application.Version)
}

func resourceToApplication(data *schema.ResourceData) *client.Application {
	application := client.Application{
		Name:    data.Get(applicationName).(string),
		Scope:   data.Get(applicationScope).(string),
		Version: data.Get(applicationVersion).(string),
	}
	application.ID = data.Id()
	return &application
}
