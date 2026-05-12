package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ---------------------------------------------------------------------------
// servicenow_rest_message
// ---------------------------------------------------------------------------

func TestAccRestMessage_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_rest_message"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccRestMessageConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_rest_message.test"),
					resource.TestCheckResourceAttr("servicenow_rest_message.test", "name", "tf-acc-rest-message"),
					resource.TestCheckResourceAttr("servicenow_rest_message.test", "rest_endpoint", "https://example.com/api"),
					resource.TestCheckResourceAttr("servicenow_rest_message.test", "description", "TF acceptance test"),
				),
			},
			{
				ResourceName:      "servicenow_rest_message.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRestMessageConfig() string {
	return `
resource "servicenow_rest_message" "test" {
  name          = "tf-acc-rest-message"
  rest_endpoint = "https://example.com/api"
  description   = "TF acceptance test"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_rest_message_header
// ---------------------------------------------------------------------------

func TestAccRestMessageHeader_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_rest_message_header"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccRestMessageHeaderConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_rest_message_header.test"),
					resource.TestCheckResourceAttr("servicenow_rest_message_header.test", "name", "Content-Type"),
					resource.TestCheckResourceAttr("servicenow_rest_message_header.test", "value", "application/json"),
				),
			},
			{
				ResourceName:      "servicenow_rest_message_header.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRestMessageHeaderConfig() string {
	return `
resource "servicenow_rest_message" "parent" {
  name          = "tf-acc-rest-msg-for-header"
  rest_endpoint = "https://example.com/api"
}

resource "servicenow_rest_message_header" "test" {
  name            = "Content-Type"
  value           = "application/json"
  rest_message_id = servicenow_rest_message.parent.id
}
`
}

// ---------------------------------------------------------------------------
// servicenow_rest_method
// ---------------------------------------------------------------------------

func TestAccRestMethod_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_rest_method"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccRestMethodConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_rest_method.test"),
					resource.TestCheckResourceAttr("servicenow_rest_method.test", "name", "get_items"),
					resource.TestCheckResourceAttr("servicenow_rest_method.test", "http_method", "get"),
				),
			},
			{
				ResourceName:      "servicenow_rest_method.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRestMethodConfig() string {
	return `
resource "servicenow_rest_message" "for_method" {
  name          = "tf-acc-rest-msg-for-method"
  rest_endpoint = "https://example.com/api"
}

resource "servicenow_rest_method" "test" {
  name            = "get_items"
  rest_message_id = servicenow_rest_message.for_method.id
  http_method     = "get"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_rest_method_header
// ---------------------------------------------------------------------------

func TestAccRestMethodHeader_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_rest_method_header"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccRestMethodHeaderConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_rest_method_header.test"),
					resource.TestCheckResourceAttr("servicenow_rest_method_header.test", "name", "Accept"),
					resource.TestCheckResourceAttr("servicenow_rest_method_header.test", "value", "application/json"),
				),
			},
			{
				ResourceName:      "servicenow_rest_method_header.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRestMethodHeaderConfig() string {
	return `
resource "servicenow_rest_message" "chain_msg" {
  name          = "tf-acc-chain-msg"
  rest_endpoint = "https://example.com/api"
}

resource "servicenow_rest_method" "chain_method" {
  name            = "chain_get"
  rest_message_id = servicenow_rest_message.chain_msg.id
  http_method     = "get"
}

resource "servicenow_rest_method_header" "test" {
  name          = "Accept"
  value         = "application/json"
  rest_method_id = servicenow_rest_method.chain_method.id
}
`
}

// ---------------------------------------------------------------------------
// servicenow_scripted_rest_api
// ---------------------------------------------------------------------------

func TestAccScriptedRestApi_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_scripted_rest_api"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccScriptedRestApiConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_scripted_rest_api.test"),
					resource.TestCheckResourceAttr("servicenow_scripted_rest_api.test", "name", "tf Acc Test API"),
					resource.TestCheckResourceAttr("servicenow_scripted_rest_api.test", "service_id", "tf_acc_test_api"),
					resource.TestCheckResourceAttr("servicenow_scripted_rest_api.test", "active", "true"),
				),
			},
			{
				ResourceName:      "servicenow_scripted_rest_api.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccScriptedRestApiConfig() string {
	return `
resource "servicenow_scripted_rest_api" "test" {
  name       = "tf Acc Test API"
  service_id = "tf_acc_test_api"
  active     = true
}
`
}

// ---------------------------------------------------------------------------
// servicenow_scripted_rest_resource
// ---------------------------------------------------------------------------

func TestAccScriptedRestResource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_scripted_rest_resource"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccScriptedRestResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_scripted_rest_resource.test"),
					resource.TestCheckResourceAttr("servicenow_scripted_rest_resource.test", "name", "tf-acc-resource"),
					resource.TestCheckResourceAttr("servicenow_scripted_rest_resource.test", "http_method", "GET"),
					resource.TestCheckResourceAttr("servicenow_scripted_rest_resource.test", "relative_path", "/items"),
					resource.TestCheckResourceAttr("servicenow_scripted_rest_resource.test", "active", "true"),
				),
			},
			{
				ResourceName:      "servicenow_scripted_rest_resource.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"operation_script"},
			},
		},
	})
}

func testAccScriptedRestResourceConfig() string {
	return `
resource "servicenow_scripted_rest_api" "parent_api" {
  name       = "tf Acc Parent API"
  service_id = "tf_acc_parent_api"
  active     = true
}

resource "servicenow_scripted_rest_resource" "test" {
  name                   = "tf-acc-resource"
  http_method            = "GET"
  relative_path          = "/items"
  web_service_definition = servicenow_scripted_rest_api.parent_api.id
  operation_script       = "(function process(/*RESTAPIRequest*/ request, /*RESTAPIResponse*/ response) {\n  response.setBody({status: 'ok'});\n})(request, response);"
  active                 = true
}
`
}

// ---------------------------------------------------------------------------
// servicenow_oauth_entity
// ---------------------------------------------------------------------------

func TestAccOAuthEntity_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_oauth_entity"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccOAuthEntityConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_oauth_entity.test"),
					resource.TestCheckResourceAttr("servicenow_oauth_entity.test", "name", "tf-acc-oauth-entity"),
					resource.TestCheckResourceAttr("servicenow_oauth_entity.test", "redirect_url", "https://example.com/callback"),
					resource.TestCheckResourceAttr("servicenow_oauth_entity.test", "access_token_lifespan", "1800"),
					resource.TestCheckResourceAttr("servicenow_oauth_entity.test", "refresh_token_lifespan", "86400"),
				),
			},
			{
				ResourceName:      "servicenow_oauth_entity.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOAuthEntityConfig() string {
	return `
resource "servicenow_oauth_entity" "test" {
  name                   = "tf-acc-oauth-entity"
  redirect_url           = "https://example.com/callback"
  access_token_lifespan  = 1800
  refresh_token_lifespan = 86400
}
`
}

// ---------------------------------------------------------------------------
// servicenow_service_catalog
// ---------------------------------------------------------------------------

func TestAccServiceCatalog_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_service_catalog"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccServiceCatalogConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_service_catalog.test"),
					resource.TestCheckResourceAttr("servicenow_service_catalog.test", "title", "TF Acc Test Catalog"),
					resource.TestCheckResourceAttr("servicenow_service_catalog.test", "description", "Created by Terraform acceptance tests"),
					resource.TestCheckResourceAttr("servicenow_service_catalog.test", "active", "true"),
				),
			},
			{
				ResourceName:      "servicenow_service_catalog.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServiceCatalogConfig() string {
	return `
resource "servicenow_service_catalog" "test" {
  title       = "TF Acc Test Catalog"
  description = "Created by Terraform acceptance tests"
  active      = true
}
`
}

// ---------------------------------------------------------------------------
// servicenow_service_catalog_category
// ---------------------------------------------------------------------------

func TestAccServiceCatalogCategory_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_service_catalog_category"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccServiceCatalogCategoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_service_catalog_category.test"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_category.test", "title", "TF Acc Test Category"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_category.test", "active", "true"),
				),
			},
			{
				ResourceName:      "servicenow_service_catalog_category.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServiceCatalogCategoryConfig() string {
	return `
resource "servicenow_service_catalog" "parent_catalog" {
  title  = "TF Acc Parent Catalog"
  active = true
}

resource "servicenow_service_catalog_category" "test" {
  title   = "TF Acc Test Category"
  catalog = servicenow_service_catalog.parent_catalog.id
  active  = true
}
`
}

// ---------------------------------------------------------------------------
// servicenow_service_catalog_item
// ---------------------------------------------------------------------------

func TestAccServiceCatalogItem_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_service_catalog_item"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccServiceCatalogItemConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_service_catalog_item.test"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_item.test", "name", "tf-acc-catalog-item"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_item.test", "short_description", "TF acceptance test item"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_item.test", "active", "true"),
				),
			},
			{
				ResourceName:      "servicenow_service_catalog_item.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServiceCatalogItemConfig() string {
	return `
resource "servicenow_service_catalog" "item_catalog" {
  title  = "TF Acc Item Catalog"
  active = true
}

resource "servicenow_service_catalog_item" "test" {
  name              = "tf-acc-catalog-item"
  short_description = "TF acceptance test item"
  service_catalogs  = servicenow_service_catalog.item_catalog.id
  active            = true
}
`
}

// ---------------------------------------------------------------------------
// servicenow_service_catalog_variable
// ---------------------------------------------------------------------------

func TestAccServiceCatalogVariable_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_service_catalog_variable"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccServiceCatalogVariableConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_service_catalog_variable.test"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_variable.test", "name", "tf_acc_var"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_variable.test", "type", "Single Line Text"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_variable.test", "order", "100"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_variable.test", "mandatory", "false"),
					resource.TestCheckResourceAttr("servicenow_service_catalog_variable.test", "active", "true"),
				),
			},
			{
				ResourceName:      "servicenow_service_catalog_variable.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServiceCatalogVariableConfig() string {
	return `
resource "servicenow_service_catalog" "var_catalog" {
  title  = "TF Acc Var Catalog"
  active = true
}

resource "servicenow_service_catalog_item" "var_item" {
  name              = "tf-acc-var-item"
  short_description = "TF acceptance test"
  service_catalogs  = servicenow_service_catalog.var_catalog.id
  active            = true
}

resource "servicenow_service_catalog_variable" "test" {
  name         = "tf_acc_var"
  question     = "Test variable question?"
  catalog_item = servicenow_service_catalog_item.var_item.id
  type         = "Single Line Text"
  order        = "100"
  mandatory    = false
  active       = true
}
`
}

// ---------------------------------------------------------------------------
// servicenow_question_choice
// ---------------------------------------------------------------------------

func TestAccQuestionChoice_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_question_choice"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccQuestionChoiceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_question_choice.test"),
					resource.TestCheckResourceAttr("servicenow_question_choice.test", "text", "Option 1"),
					resource.TestCheckResourceAttr("servicenow_question_choice.test", "value", "option_1"),
					resource.TestCheckResourceAttr("servicenow_question_choice.test", "order", "100"),
				),
			},
			{
				ResourceName:      "servicenow_question_choice.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccQuestionChoiceConfig() string {
	return `
resource "servicenow_service_catalog" "qc_catalog" {
  title  = "TF Acc QC Catalog"
  active = true
}

resource "servicenow_service_catalog_item" "qc_item" {
  name              = "tf-acc-qc-item"
  short_description = "TF acceptance test"
  service_catalogs  = servicenow_service_catalog.qc_catalog.id
  active            = true
}

resource "servicenow_service_catalog_variable" "qc_var" {
  name         = "tf_acc_qc_var"
  question     = "Pick an option?"
  catalog_item = servicenow_service_catalog_item.qc_item.id
  type         = "Select Box"
  order        = "100"
  mandatory    = false
  active       = true
}

resource "servicenow_question_choice" "test" {
  text     = "Option 1"
  value    = "option_1"
  question = servicenow_service_catalog_variable.qc_var.id
  order    = "100"
}
`
}

// ---------------------------------------------------------------------------
// servicenow_server
// ---------------------------------------------------------------------------

func TestAccServer_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      checkDestroy("servicenow_server"),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + testAccServerConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_server.test"),
					resource.TestCheckResourceAttr("servicenow_server.test", "name", "tf-acc-test-server"),
					resource.TestCheckResourceAttr("servicenow_server.test", "ip_address", "192.168.100.100"),
				),
			},
			{
				ResourceName:      "servicenow_server.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServerConfig() string {
	return `
resource "servicenow_server" "test" {
  name       = "tf-acc-test-server"
  ip_address = "192.168.100.100"
}
`
}
