package servicenow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/resources"
)

// Provider is a Terraform Provider that manages objects in a ServiceNow instance.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"instance_url": {
				Type:         schema.TypeString,
				Description:  "The URL of the ServiceNow instance to work with. May also be set via the SERVICENOW_INSTANCE_URL environment variable.",
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SERVICENOW_INSTANCE_URL", nil),
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username used to manage resources via Basic authentication. May also be set via the SERVICENOW_USERNAME environment variable.",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICENOW_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password of the user to manage resources. May also be set via the SERVICENOW_PASSWORD environment variable.",
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICENOW_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"servicenow_alias":                      resources.ResourceAlias(),
			"servicenow_application":                resources.ResourceApplication(),
			"servicenow_application_menu":           resources.ResourceApplicationMenu(),
			"servicenow_application_module":         resources.ResourceApplicationModule(),
			"servicenow_basic_auth_credential":      resources.ResourceBasicAuthCredential(),
			"servicenow_content_css":                resources.ResourceContentCSS(),
			"servicenow_css_include":                resources.ResourceCSSInclude(),
			"servicenow_css_include_relation":       resources.ResourceCSSIncludeRelation(),
			"servicenow_db_table":                   resources.ResourceDBTable(),
			"servicenow_extension_point":            resources.ResourceExtensionPoint(),
			"servicenow_http_connection":            resources.ResourceHttpConnection(),
			"servicenow_js_include":                 resources.ResourceJsInclude(),
			"servicenow_js_include_relation":        resources.ResourceJsIncludeRelation(),
			"servicenow_oauth_entity":               resources.ResourceOAuthEntity(),
			"servicenow_question_choice":            resources.ResourceQuestionChoice(),
			"servicenow_role":                       resources.ResourceRole(),
			"servicenow_rest_message":               resources.ResourceRestMessage(),
			"servicenow_rest_message_header":        resources.ResourceRestMessageHeader(),
			"servicenow_rest_method":                resources.ResourceRestMethod(),
			"servicenow_rest_method_header":         resources.ResourceRestMethodHeader(),
			"servicenow_scripted_rest_api":          resources.ResourceScriptedRestApi(),
			"servicenow_scripted_rest_resource":     resources.ResourceScriptedRestResource(),
			"servicenow_script_include":             resources.ResourceScriptInclude(),
			"servicenow_server":                     resources.ResourceServer(),
			"servicenow_service_catalog":            resources.ResourceServiceCatalog(),
			"servicenow_service_catalog_category":   resources.ResourceServiceCatalogCategory(),
			"servicenow_service_catalog_item":       resources.ResourceServiceCatalogItem(),
			"servicenow_service_catalog_variable":   resources.ResourceServiceCatalogVariable(),
			"servicenow_system_property":            resources.ResourceSystemProperty(),
			"servicenow_system_property_category":   resources.ResourceSystemPropertyCategory(),
			"servicenow_system_property_relation":   resources.ResourceSystemPropertyRelation(),
			"servicenow_ui_macro":                   resources.ResourceUIMacro(),
			"servicenow_ui_page":                    resources.ResourceUIPage(),
			"servicenow_ui_script":                  resources.ResourceUIScript(),
			"servicenow_widget":                     resources.ResourceWidget(),
			"servicenow_widget_dependency":          resources.ResourceWidgetDependency(),
			"servicenow_widget_dependency_relation": resources.ResourceWidgetDependencyRelation(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"servicenow_acl":                      resources.DataSourceACL(),
			"servicenow_application":              resources.DataSourceApplication(),
			"servicenow_application_category":     resources.DataSourceApplicationCategory(),
			"servicenow_db_table":                 resources.DataSourceDBTable(),
			"servicenow_role":                     resources.DataSourceRole(),
			"servicenow_service_catalog":          resources.DataSourceServiceCatalog(),
			"servicenow_service_catalog_category": resources.DataSourceServiceCatalogCategory(),
			"servicenow_system_property":          resources.DataSourceSystemProperty(),
			"servicenow_system_property_category": resources.DataSourceSystemPropertyCategory(),
		},
		ConfigureContextFunc: configure,
	}
}

func configure(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	c := client.NewClient(
		data.Get("instance_url").(string),
		data.Get("username").(string),
		data.Get("password").(string),
	)
	c.UserAgent = "terraform-provider-servicenow/" + providerVersion
	return c, nil
}

// providerVersion is set from main via SetVersion at startup.
var providerVersion = "dev"

// SetVersion sets the provider version string used in User-Agent headers.
func SetVersion(v string) {
	if v != "" {
		providerVersion = v
	}
}

// Version returns the current provider version.
func Version() string {
	return providerVersion
}
