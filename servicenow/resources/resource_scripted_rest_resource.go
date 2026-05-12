package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const scriptedRestResourceName = "name"
const scriptedRestResourceActive = "active"
const scriptedRestResourceEnforceACL = "enforce_acl"
const scriptedRestResourceRequiresACLAuthorization = "requires_acl_authorization"
const scriptedRestResourceRequiresAuthentication = "requires_authentication"
const scriptedRestResourceRequiresSNCInternalRole = "requires_snc_internal_role"
const scriptedRestResourceProduces = "produces"
const scriptedRestResourceShortDescription = "short_description"
const scriptedRestResourceOperationScript = "operation_script"
const scriptedRestResourceRelativePath = "relative_path"
const scriptedRestResourceRequestExample = "request_example"
const scriptedRestResourceHTTPMethod = "http_method"
const scriptedRestResourceConsumes = "consumes"
const scriptedRestResourceOperationURI = "operation_uri"
const scriptedRestResourceWebServiceDefinition = "web_service_definition"
const scriptedRestResourceWebServiceVersion = "web_service_version"

// ResourceScriptedRestResource manages a Scripted Rest Resource in ServiceNow.
func ResourceScriptedRestResource() *schema.Resource {

	return &schema.Resource{
		Description: "`servicenow_scripted_rest_resource` manages a scripted REST resource within ServiceNow.",

		CreateContext: createResourceScriptedRestResource,
		ReadContext:   readResourceScriptedRestResource,
		UpdateContext: updateResourceScriptedRestResource,
		DeleteContext: deleteResourceScriptedRestResource,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			scriptedRestResourceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the API resource. Appears in API documentation.",
			},
			scriptedRestResourceHTTPMethod: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The HTTP method that maps to this record. Can be 'GET', 'POST', 'PUT', 'PATCH' or 'DELETE'.",
				ValidateFunc: validation.StringInSlice([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}, false),
			},
			scriptedRestResourceOperationScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The script that implements the resource.",
			},
			scriptedRestResourceWebServiceDefinition: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parent API this resource belongs to.",
			},
			scriptedRestResourceActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Activates the resource. Inactive resources cannot serve requests.",
			},
			scriptedRestResourceEnforceACL: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ACLs to enforce when accessing resources. Individual resources may override this value.",
			},
			scriptedRestResourceRequiresACLAuthorization: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enforce ACLs when this resource is accessed.",
			},
			scriptedRestResourceRequiresAuthentication: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether requests must be authenticated to access this resource.",
			},
			scriptedRestResourceRequiresSNCInternalRole: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether requests must be authenticated with SNC Internal Role to access this resource.",
			},
			scriptedRestResourceConsumes: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "application/json",
				Description: "Default supported request formats.",
			},
			scriptedRestResourceProduces: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "application/json",
				Description: "Default supported response formats.",
			},
			scriptedRestResourceShortDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the API. Appears in API documentation.",
			},
			scriptedRestResourceRelativePath: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path of this resource relative to the base API path. Can contain templatized path paramenters such as /{id}.",
			},
			scriptedRestResourceRequestExample: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An example of a request sent to this resource.",
			},
			scriptedRestResourceOperationURI: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resolved path of this resource including base API path, version, and relative path.",
			},
			scriptedRestResourceWebServiceVersion: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The version of the parent API this resource belongs to.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceScriptedRestResource(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scriptedRestResource := &client.ScriptedRestResource{}
	if err := snowClient.GetObject(ctx, client.EndpointScriptedRestResource, data.Id(), scriptedRestResource); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromScriptedRestResource(data, scriptedRestResource)

	return nil
}

func createResourceScriptedRestResource(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	scriptedRestResource := resourceToScriptedRestResource(data)
	if err := snowClient.CreateObject(ctx, client.EndpointScriptedRestResource, scriptedRestResource); err != nil {
		return diag.FromErr(err)
	}

	resourceFromScriptedRestResource(data, scriptedRestResource)

	return readResourceScriptedRestResource(ctx, data, serviceNowClient)
}

func updateResourceScriptedRestResource(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointScriptedRestResource, resourceToScriptedRestResource(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceScriptedRestResource(ctx, data, serviceNowClient)
}

func deleteResourceScriptedRestResource(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointScriptedRestResource, data.Id()))
}

func resourceFromScriptedRestResource(data *schema.ResourceData, scriptedRestResource *client.ScriptedRestResource) {
	data.SetId(scriptedRestResource.ID)
	data.Set(scriptedRestResourceName, scriptedRestResource.Name)
	data.Set(scriptedRestResourceActive, scriptedRestResource.Active)
	data.Set(scriptedRestResourceEnforceACL, scriptedRestResource.EnforceACL)
	data.Set(scriptedRestResourceRequiresACLAuthorization, scriptedRestResource.RequiresACLAuthorization)
	data.Set(scriptedRestResourceRequiresAuthentication, scriptedRestResource.RequiresAuthentication)
	data.Set(scriptedRestResourceRequiresSNCInternalRole, scriptedRestResource.RequiresSNCInternalRole)
	data.Set(scriptedRestResourceProduces, scriptedRestResource.Produces)
	data.Set(scriptedRestResourceShortDescription, scriptedRestResource.ShortDescription)
	data.Set(scriptedRestResourceOperationScript, scriptedRestResource.OperationScript)
	data.Set(scriptedRestResourceRelativePath, scriptedRestResource.RelativePath)
	data.Set(scriptedRestResourceRequestExample, scriptedRestResource.RequestExample)
	data.Set(scriptedRestResourceHTTPMethod, scriptedRestResource.HTTPMethod)
	data.Set(scriptedRestResourceConsumes, scriptedRestResource.Consumes)
	data.Set(scriptedRestResourceWebServiceDefinition, scriptedRestResource.WebServiceDefinition)
	data.Set(scriptedRestResourceWebServiceVersion, scriptedRestResource.WebServiceVersion)
	data.Set(commonProtectionPolicy, scriptedRestResource.ProtectionPolicy)
	data.Set(commonScope, scriptedRestResource.Scope)
}

func resourceToScriptedRestResource(data *schema.ResourceData) *client.ScriptedRestResource {
	scriptedRestResource := client.ScriptedRestResource{
		Name:                     data.Get(scriptedRestResourceName).(string),
		Active:                   data.Get(scriptedRestResourceActive).(bool),
		EnforceACL:               data.Get(scriptedRestResourceEnforceACL).(string),
		RequiresACLAuthorization: data.Get(scriptedRestResourceRequiresACLAuthorization).(bool),
		RequiresAuthentication:   data.Get(scriptedRestResourceRequiresAuthentication).(bool),
		RequiresSNCInternalRole:  data.Get(scriptedRestResourceRequiresSNCInternalRole).(bool),
		ShortDescription:         data.Get(scriptedRestResourceShortDescription).(string),
		OperationScript:          data.Get(scriptedRestResourceOperationScript).(string),
		RelativePath:             data.Get(scriptedRestResourceRelativePath).(string),
		RequestExample:           data.Get(scriptedRestResourceRequestExample).(string),
		HTTPMethod:               data.Get(scriptedRestResourceHTTPMethod).(string),
		WebServiceDefinition:     data.Get(scriptedRestResourceWebServiceDefinition).(string),
		WebServiceVersion:        data.Get(scriptedRestResourceWebServiceVersion).(string),
		Produces:                 data.Get(scriptedRestResourceProduces).(string),
		Consumes:                 data.Get(scriptedRestResourceConsumes).(string),
		ProducesCustomized:       true,
		ConsumesCustomized:       true,
	}

	scriptedRestResource.ID = data.Id()
	scriptedRestResource.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	scriptedRestResource.Scope = data.Get(commonScope).(string)
	return &scriptedRestResource
}
