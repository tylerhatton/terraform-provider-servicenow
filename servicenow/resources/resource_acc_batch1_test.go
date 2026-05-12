package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ---------------------------------------------------------------------------
// servicenow_alias
// ---------------------------------------------------------------------------

func TestAccAlias_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_alias"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccAliasConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_alias.test"),
					resource.TestCheckResourceAttr("servicenow_alias.test", "name", "tf-acc-alias"),
					resource.TestCheckResourceAttr("servicenow_alias.test", "type", "credential"),
				),
			},
			{
				ResourceName:      "servicenow_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliasConfig() string {
	return `
resource "servicenow_alias" "test" {
  name = "tf-acc-alias"
  type = "credential"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_application
// ---------------------------------------------------------------------------

func TestAccApplication_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_application"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccApplicationConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_application.test"),
					resource.TestCheckResourceAttr("servicenow_application.test", "name", "TF Acc Test App"),
					resource.TestCheckResourceAttr("servicenow_application.test", "scope", "x_tfacc_test"),
					resource.TestCheckResourceAttr("servicenow_application.test", "version", "1.0.0"),
				),
			},
			{
				ResourceName:      "servicenow_application.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApplicationConfig() string {
	return `
resource "servicenow_application" "test" {
  name    = "TF Acc Test App"
  scope   = "x_tfacc_test"
  version = "1.0.0"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_basic_auth_credential
// ---------------------------------------------------------------------------

func TestAccBasicAuthCredential_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_basic_auth_credential"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccBasicAuthCredentialConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_basic_auth_credential.test"),
					resource.TestCheckResourceAttr("servicenow_basic_auth_credential.test", "name", "tf-acc-basic-auth"),
					resource.TestCheckResourceAttr("servicenow_basic_auth_credential.test", "username", "testuser"),
					resource.TestCheckResourceAttrSet("servicenow_basic_auth_credential.test", "credential_alias"),
				),
			},
			{
				ResourceName:            "servicenow_basic_auth_credential.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccBasicAuthCredentialConfig() string {
	return `
resource "servicenow_alias" "cred_alias" {
  name = "tf-acc-basic-auth-alias"
  type = "credential"
}

resource "servicenow_basic_auth_credential" "test" {
  name             = "tf-acc-basic-auth"
  username         = "testuser"
  password         = "testpass123"
  credential_alias = servicenow_alias.cred_alias.id
}
`
}

// ---------------------------------------------------------------------------
// servicenow_role
// ---------------------------------------------------------------------------

func TestAccRole_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_role"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccRoleConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_role.test"),
					resource.TestCheckResourceAttr("servicenow_role.test", "suffix", "tf_acc_test_role"),
				),
			},
			{
				ResourceName:      "servicenow_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRoleConfig() string {
	return `
resource "servicenow_role" "test" {
  suffix = "tf_acc_test_role"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_script_include
// ---------------------------------------------------------------------------

func TestAccScriptInclude_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_script_include"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccScriptIncludeConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_script_include.test"),
					resource.TestCheckResourceAttr("servicenow_script_include.test", "name", "TfAccTestScriptInclude"),
					resource.TestCheckResourceAttr("servicenow_script_include.test", "active", "true"),
				),
			},
			{
				ResourceName:      "servicenow_script_include.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccScriptIncludeConfig() string {
	return `
resource "servicenow_script_include" "test" {
  name   = "TfAccTestScriptInclude"
  script = "var TfAccTestScriptInclude = Class.create();\nTfAccTestScriptInclude.prototype = {\n  initialize: function() {},\n  type: 'TfAccTestScriptInclude'\n};"
  active = true
}
`
}

// ---------------------------------------------------------------------------
// servicenow_ui_macro
// ---------------------------------------------------------------------------

func TestAccUIMacro_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_ui_macro"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccUIMacroConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_ui_macro.test"),
					resource.TestCheckResourceAttr("servicenow_ui_macro.test", "name", "tf_acc_test_macro"),
				),
			},
			{
				ResourceName:      "servicenow_ui_macro.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUIMacroConfig() string {
	return `
resource "servicenow_ui_macro" "test" {
  name = "tf_acc_test_macro"
  xml  = "<j:jelly xmlns:j='jelly:core' xmlns:g='glide' xmlns:j2='null' xmlns:g2='null'><h1>Test Macro</h1></j:jelly>"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_ui_page
// ---------------------------------------------------------------------------

func TestAccUIPage_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_ui_page"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccUIPageConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_ui_page.test"),
					resource.TestCheckResourceAttr("servicenow_ui_page.test", "name", "tf_acc_test_page"),
					resource.TestCheckResourceAttr("servicenow_ui_page.test", "category", "general"),
				),
			},
			{
				ResourceName:      "servicenow_ui_page.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUIPageConfig() string {
	return `
resource "servicenow_ui_page" "test" {
  name              = "tf_acc_test_page"
  html              = "<h1>Test Page</h1>"
  client_script     = "// client script"
  processing_script = "// processing script"
  category          = "general"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_ui_script
// ---------------------------------------------------------------------------

func TestAccUIScript_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_ui_script"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccUIScriptConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_ui_script.test"),
					resource.TestCheckResourceAttr("servicenow_ui_script.test", "name", "tf_acc_test_ui_script"),
					resource.TestCheckResourceAttr("servicenow_ui_script.test", "active", "true"),
				),
			},
			{
				ResourceName:      "servicenow_ui_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUIScriptConfig() string {
	return `
resource "servicenow_ui_script" "test" {
  name   = "tf_acc_test_ui_script"
  script = "// TF acceptance test script\nconsole.log('hello');"
  active = true
}
`
}

// ---------------------------------------------------------------------------
// servicenow_system_property
// ---------------------------------------------------------------------------

func TestAccSystemProperty_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_system_property"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccSystemPropertyConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_system_property.test"),
					resource.TestCheckResourceAttr("servicenow_system_property.test", "suffix", "tf.acc.test.property"),
					resource.TestCheckResourceAttr("servicenow_system_property.test", "type", "string"),
				),
			},
			{
				ResourceName:      "servicenow_system_property.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSystemPropertyConfig() string {
	return `
resource "servicenow_system_property" "test" {
  suffix = "tf.acc.test.property"
  type   = "string"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_system_property_category
// ---------------------------------------------------------------------------

func TestAccSystemPropertyCategory_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_system_property_category"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccSystemPropertyCategoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_system_property_category.test"),
					resource.TestCheckResourceAttr("servicenow_system_property_category.test", "name", "TF Acc Test Category"),
				),
			},
			{
				ResourceName:      "servicenow_system_property_category.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSystemPropertyCategoryConfig() string {
	return `
resource "servicenow_system_property_category" "test" {
  name = "TF Acc Test Category"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_extension_point
// ---------------------------------------------------------------------------

func TestAccExtensionPoint_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_extension_point"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccExtensionPointConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_extension_point.test"),
					resource.TestCheckResourceAttr("servicenow_extension_point.test", "name", "TF Acc Test Extension Point"),
					resource.TestCheckResourceAttr("servicenow_extension_point.test", "description", "Created by Terraform acceptance tests"),
				),
			},
			{
				ResourceName:      "servicenow_extension_point.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccExtensionPointConfig() string {
	return `
resource "servicenow_extension_point" "test" {
  name        = "TF Acc Test Extension Point"
  description = "Created by Terraform acceptance tests"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_http_connection
// ---------------------------------------------------------------------------

func TestAccHttpConnection_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_http_connection"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccHttpConnectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_http_connection.test"),
					resource.TestCheckResourceAttr("servicenow_http_connection.test", "name", "tf-acc-http-connection"),
					resource.TestCheckResourceAttr("servicenow_http_connection.test", "connection_url", "https://example.com"),
					resource.TestCheckResourceAttr("servicenow_http_connection.test", "active", "true"),
					resource.TestCheckResourceAttrSet("servicenow_http_connection.test", "connection_alias"),
				),
			},
			{
				ResourceName:      "servicenow_http_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccHttpConnectionConfig() string {
	return `
resource "servicenow_alias" "conn_alias" {
  name            = "tf-acc-http-conn-alias"
  type            = "connection"
  connection_type = "http_connection"
}

resource "servicenow_http_connection" "test" {
  name             = "tf-acc-http-connection"
  connection_alias = servicenow_alias.conn_alias.id
  connection_url   = "https://example.com"
  active           = true
}
`
}
