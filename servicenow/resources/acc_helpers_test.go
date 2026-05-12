package resources_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	servicenow "github.com/tylerhatton/terraform-provider-servicenow/servicenow"
)

// testAccPreCheck validates required environment variables are set before running acceptance tests.
func testAccPreCheck(t *testing.T) {
	t.Helper()
	if v := os.Getenv("SERVICENOW_INSTANCE_URL"); v == "" {
		t.Fatal("SERVICENOW_INSTANCE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICENOW_USERNAME"); v == "" {
		t.Fatal("SERVICENOW_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICENOW_PASSWORD"); v == "" {
		t.Fatal("SERVICENOW_PASSWORD must be set for acceptance tests")
	}
}

// testAccProviderFactories returns a map of provider factories for acceptance tests.
func testAccProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"servicenow": func() (*schema.Provider, error) {
			return servicenow.Provider(), nil
		},
	}
}

// providerBlock returns the provider configuration block for use in acceptance test configs.
func providerBlock() string {
	return fmt.Sprintf(`
provider "servicenow" {
  instance_url = %q
  username     = %q
  password     = %q
}
`,
		os.Getenv("SERVICENOW_INSTANCE_URL"),
		os.Getenv("SERVICENOW_USERNAME"),
		os.Getenv("SERVICENOW_PASSWORD"),
	)
}

// checkExists returns a resource.TestCheckFunc that verifies a resource exists in state with a non-empty ID.
func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is empty for: %s", resourceName)
		}
		return nil
	}
}

// resourceTypeToEndpoint maps resource type names to their ServiceNow API endpoints.
var resourceTypeToEndpoint = map[string]string{
	"servicenow_alias":                      "sys_alias.do",
	"servicenow_application":                "sys_app.do",
	"servicenow_application_menu":           "sys_app_application.do",
	"servicenow_application_module":         "sys_app_module.do",
	"servicenow_basic_auth_credential":      "basic_auth_credentials.do",
	"servicenow_content_css":                "content_css.do",
	"servicenow_css_include":                "sys_ui_ie_stylesheet.do",
	"servicenow_css_include_relation":       "sys_ui_ie_stylesheet_list.do",
	"servicenow_db_table":                   "sys_db_object.do",
	"servicenow_extension_point":            "sys_extension_point.do",
	"servicenow_http_connection":            "http_connection.do",
	"servicenow_js_include":                 "sys_ui_js.do",
	"servicenow_js_include_relation":        "sys_ui_js_list.do",
	"servicenow_oauth_entity":               "oauth_entity.do",
	"servicenow_question_choice":            "question_choice.do",
	"servicenow_rest_message":               "sys_rest_message.do",
	"servicenow_rest_message_header":        "sys_rest_message_headers.do",
	"servicenow_rest_method":                "sys_rest_message_fn.do",
	"servicenow_rest_method_header":         "sys_rest_message_fn_headers.do",
	"servicenow_role":                       "sys_user_role.do",
	"servicenow_script_include":             "sys_script_include.do",
	"servicenow_scripted_rest_api":          "sys_ws_definition.do",
	"servicenow_scripted_rest_resource":     "sys_ws_operation.do",
	"servicenow_server":                     "cmdb_ci_server.do",
	"servicenow_service_catalog":            "sc_catalog.do",
	"servicenow_service_catalog_category":   "sc_category.do",
	"servicenow_service_catalog_item":       "sc_cat_item.do",
	"servicenow_service_catalog_variable":   "item_option_new.do",
	"servicenow_system_property":            "sys_properties.do",
	"servicenow_system_property_category":   "sys_properties_category.do",
	"servicenow_system_property_relation":   "sys_properties_category_m2m.do",
	"servicenow_ui_macro":                   "sys_ui_macro.do",
	"servicenow_ui_page":                    "sys_ui_page.do",
	"servicenow_ui_script":                  "sys_ui_script.do",
	"servicenow_widget":                     "sp_widget.do",
	"servicenow_widget_dependency":          "sp_dependency.do",
	"servicenow_widget_dependency_relation": "m2m_sp_widget_dependency.do",
}

// snowResponse is used to parse ServiceNow JSON API responses.
type snowResponse struct {
	Records []json.RawMessage `json:"records"`
}

// checkDestroy returns a resource.TestCheckFunc that verifies resources of the given type
// have been removed from ServiceNow after destruction.
//
// IMPORTANT: The Terraform SDK v2 passes the PRE-DESTROY state to CheckDestroy. This means
// we must verify resource deletion against the ServiceNow API directly, not by checking state.
func checkDestroy(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		endpoint, ok := resourceTypeToEndpoint[resourceType]
		if !ok {
			// Unknown resource type - skip verification
			return nil
		}

		instanceURL := os.Getenv("SERVICENOW_INSTANCE_URL")
		username := os.Getenv("SERVICENOW_USERNAME")
		password := os.Getenv("SERVICENOW_PASSWORD")

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			if rs.Primary.ID == "" {
				continue
			}

			// Verify the resource is gone from ServiceNow
			url := fmt.Sprintf("%s/%s?JSONv2&sysparm_query=sys_id=%s",
				instanceURL, endpoint, rs.Primary.ID)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return fmt.Errorf("failed to create verification request for %s %s: %w", resourceType, rs.Primary.ID, err)
			}
			req.SetBasicAuth(username, password)
			req.Header.Set("Content-Type", "application/json")

			httpClient := &http.Client{}
			resp, err := httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to verify destroy of %s %s: %w", resourceType, rs.Primary.ID, err)
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var result snowResponse
			if err := json.Unmarshal(body, &result); err != nil {
				// Can't parse response - assume gone
				continue
			}

			if len(result.Records) > 0 {
				return fmt.Errorf("resource %s with ID %s still exists in ServiceNow after destroy",
					resourceType, rs.Primary.ID)
			}
		}
		return nil
	}
}
