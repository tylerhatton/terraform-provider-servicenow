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
			"servicenow_acl":                        resources.ResourceACL(),
			"servicenow_alias":                      resources.ResourceAlias(),
			"servicenow_application":                resources.ResourceApplication(),
			"servicenow_application_menu":           resources.ResourceApplicationMenu(),
			"servicenow_application_module":         resources.ResourceApplicationModule(),
			"servicenow_assignment_rule":            resources.ResourceAssignmentRule(),
			"servicenow_basic_auth_credential":      resources.ResourceBasicAuthCredential(),
			"servicenow_business_rule":              resources.ResourceBusinessRule(),
			"servicenow_certificate":                resources.ResourceCertificate(),
			"servicenow_choice":                     resources.ResourceChoice(),
			"servicenow_client_script":              resources.ResourceClientScript(),
			"servicenow_content_css":                resources.ResourceContentCSS(),
			"servicenow_css_include":                resources.ResourceCSSInclude(),
			"servicenow_css_include_relation":       resources.ResourceCSSIncludeRelation(),
			"servicenow_data_lookup":                resources.ResourceDataLookup(),
			"servicenow_db_table":                   resources.ResourceDBTable(),
			"servicenow_dictionary":                 resources.ResourceDictionary(),
			"servicenow_email_template":             resources.ResourceEmailTemplate(),
			"servicenow_encryption_context":         resources.ResourceEncryptionContext(),
			"servicenow_extension_point":            resources.ResourceExtensionPoint(),
			"servicenow_flow":                       resources.ResourceFlow(),
			"servicenow_group":                      resources.ResourceGroup(),
			"servicenow_group_member":               resources.ResourceGroupMember(),
			"servicenow_group_role":                 resources.ResourceGroupRole(),
			"servicenow_http_connection":            resources.ResourceHttpConnection(),
			"servicenow_jdbc_connection":            resources.ResourceJdbcConnection(),
			"servicenow_js_include":                 resources.ResourceJsInclude(),
			"servicenow_js_include_relation":        resources.ResourceJsIncludeRelation(),
			"servicenow_mid_server":                 resources.ResourceMidServer(),
			"servicenow_notification":               resources.ResourceNotification(),
			"servicenow_oauth_entity":               resources.ResourceOAuthEntity(),
			"servicenow_question_choice":            resources.ResourceQuestionChoice(),
			"servicenow_role":                       resources.ResourceRole(),
			"servicenow_rest_message":               resources.ResourceRestMessage(),
			"servicenow_rest_message_header":        resources.ResourceRestMessageHeader(),
			"servicenow_rest_method":                resources.ResourceRestMethod(),
			"servicenow_rest_method_header":         resources.ResourceRestMethodHeader(),
			"servicenow_scheduled_job":              resources.ResourceScheduledJob(),
			"servicenow_script_action":              resources.ResourceScriptAction(),
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
			"servicenow_transform_entry":            resources.ResourceTransformEntry(),
			"servicenow_transform_map":              resources.ResourceTransformMap(),
			"servicenow_ui_action":                  resources.ResourceUIAction(),
			"servicenow_ui_macro":                   resources.ResourceUIMacro(),
			"servicenow_ui_page":                    resources.ResourceUIPage(),
			"servicenow_ui_policy":                  resources.ResourceUIPolicy(),
			"servicenow_ui_policy_action":           resources.ResourceUIPolicyAction(),
			"servicenow_ui_script":                  resources.ResourceUIScript(),
			"servicenow_user":                       resources.ResourceUser(),
			"servicenow_user_role":                  resources.ResourceUserRole(),
			"servicenow_widget":                     resources.ResourceWidget(),
			"servicenow_widget_dependency":          resources.ResourceWidgetDependency(),
			"servicenow_widget_dependency_relation": resources.ResourceWidgetDependencyRelation(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"servicenow_acl":                      resources.DataSourceACL(),
			"servicenow_alias":                    resources.DataSourceAlias(),
			"servicenow_application":              resources.DataSourceApplication(),
			"servicenow_application_category":     resources.DataSourceApplicationCategory(),
			"servicenow_application_menu":         resources.DataSourceApplicationMenu(),
			"servicenow_application_module":       resources.DataSourceApplicationModule(),
			"servicenow_assignment_rule":          resources.DataSourceAssignmentRule(),
			"servicenow_basic_auth_credential":    resources.DataSourceBasicAuthCredential(),
			"servicenow_business_rule":            resources.DataSourceBusinessRule(),
			"servicenow_certificate":              resources.DataSourceCertificate(),
			"servicenow_choice":                   resources.DataSourceChoice(),
			"servicenow_client_script":            resources.DataSourceClientScript(),
			"servicenow_content_css":              resources.DataSourceContentCSS(),
			"servicenow_css_include":              resources.DataSourceCSSInclude(),
			"servicenow_data_lookup":              resources.DataSourceDataLookup(),
			"servicenow_db_table":                 resources.DataSourceDBTable(),
			"servicenow_dictionary":               resources.DataSourceDictionary(),
			"servicenow_email_template":           resources.DataSourceEmailTemplate(),
			"servicenow_encryption_context":       resources.DataSourceEncryptionContext(),
			"servicenow_extension_point":          resources.DataSourceExtensionPoint(),
			"servicenow_flow":                     resources.DataSourceFlow(),
			"servicenow_group":                    resources.DataSourceGroup(),
			"servicenow_http_connection":          resources.DataSourceHttpConnection(),
			"servicenow_jdbc_connection":          resources.DataSourceJdbcConnection(),
			"servicenow_js_include":               resources.DataSourceJsInclude(),
			"servicenow_mid_server":               resources.DataSourceMidServer(),
			"servicenow_notification":             resources.DataSourceNotification(),
			"servicenow_oauth_entity":             resources.DataSourceOAuthEntity(),
			"servicenow_rest_message":             resources.DataSourceRestMessage(),
			"servicenow_role":                     resources.DataSourceRole(),
			"servicenow_scheduled_job":            resources.DataSourceScheduledJob(),
			"servicenow_script_action":            resources.DataSourceScriptAction(),
			"servicenow_script_include":           resources.DataSourceScriptInclude(),
			"servicenow_scripted_rest_api":        resources.DataSourceScriptedRestApi(),
			"servicenow_server":                   resources.DataSourceServer(),
			"servicenow_service_catalog":          resources.DataSourceServiceCatalog(),
			"servicenow_service_catalog_category": resources.DataSourceServiceCatalogCategory(),
			"servicenow_service_catalog_item":     resources.DataSourceServiceCatalogItem(),
			"servicenow_system_property":          resources.DataSourceSystemProperty(),
			"servicenow_system_property_category": resources.DataSourceSystemPropertyCategory(),
			"servicenow_transform_map":            resources.DataSourceTransformMap(),
			"servicenow_ui_action":                resources.DataSourceUIAction(),
			"servicenow_ui_macro":                 resources.DataSourceUIMacro(),
			"servicenow_ui_page":                  resources.DataSourceUIPage(),
			"servicenow_ui_policy":                resources.DataSourceUIPolicy(),
			"servicenow_ui_script":                resources.DataSourceUIScript(),
			"servicenow_user":                     resources.DataSourceUser(),
			"servicenow_widget":                   resources.DataSourceWidget(),
			"servicenow_widget_dependency":        resources.DataSourceWidgetDependency(),
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
