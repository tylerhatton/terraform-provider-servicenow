package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ---------------------------------------------------------------------------
// Resource: servicenow_widget
// Requires Service Portal plugin; skip gracefully when not available.
// ---------------------------------------------------------------------------

func TestAccResourceWidget_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SERVICE_PORTAL") == "" {
		t.Skip("SERVICENOW_SERVICE_PORTAL must be set to run Service Portal acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_widget" "test" {
  identifier  = "tf-acc-widget"
  name        = "TF Acc Test Widget"
  template    = "<div>Test</div>"
  description = "TF acceptance test widget"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_widget.test"),
					resource.TestCheckResourceAttr("servicenow_widget.test", "name", "TF Acc Test Widget"),
					resource.TestCheckResourceAttr("servicenow_widget.test", "identifier", "tf-acc-widget"),
					resource.TestCheckResourceAttrSet("servicenow_widget.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_widget.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_widget_dependency
// ---------------------------------------------------------------------------

func TestAccResourceWidgetDependency_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SERVICE_PORTAL") == "" {
		t.Skip("SERVICENOW_SERVICE_PORTAL must be set to run Service Portal acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_widget_dependency" "test" {
  name = "tf-acc-widget-dep"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_widget_dependency.test"),
					resource.TestCheckResourceAttr("servicenow_widget_dependency.test", "name", "tf-acc-widget-dep"),
					resource.TestCheckResourceAttrSet("servicenow_widget_dependency.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_widget_dependency.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_widget_dependency_relation
// ---------------------------------------------------------------------------

func TestAccResourceWidgetDependencyRelation_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SERVICE_PORTAL") == "" {
		t.Skip("SERVICENOW_SERVICE_PORTAL must be set to run Service Portal acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_widget" "for_rel" {
  identifier = "tf-acc-widget-for-rel"
  name       = "TF Acc Widget For Relation"
  template   = "<div></div>"
}
resource "servicenow_widget_dependency" "for_rel" {
  name = "tf-acc-dep-for-rel"
}
resource "servicenow_widget_dependency_relation" "test" {
  widget_id     = servicenow_widget.for_rel.id
  dependency_id = servicenow_widget_dependency.for_rel.id
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_widget_dependency_relation.test"),
					resource.TestCheckResourceAttrSet("servicenow_widget_dependency_relation.test", "widget_id"),
					resource.TestCheckResourceAttrSet("servicenow_widget_dependency_relation.test", "dependency_id"),
					resource.TestCheckResourceAttrSet("servicenow_widget_dependency_relation.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_widget_dependency_relation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_css_include
// ---------------------------------------------------------------------------

func TestAccResourceCSSInclude_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SERVICE_PORTAL") == "" {
		t.Skip("SERVICENOW_SERVICE_PORTAL must be set to run Service Portal acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_css_include" "test" {
  name   = "tf-acc-css-include"
  source = "url"
  url    = "https://example.com/style.css"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_css_include.test"),
					resource.TestCheckResourceAttr("servicenow_css_include.test", "name", "tf-acc-css-include"),
					resource.TestCheckResourceAttr("servicenow_css_include.test", "url", "https://example.com/style.css"),
					resource.TestCheckResourceAttrSet("servicenow_css_include.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_css_include.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_css_include_relation
// ---------------------------------------------------------------------------

func TestAccResourceCSSIncludeRelation_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SERVICE_PORTAL") == "" {
		t.Skip("SERVICENOW_SERVICE_PORTAL must be set to run Service Portal acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_css_include" "for_rel" {
  name   = "tf-acc-css-for-rel"
  source = "url"
  url    = "https://example.com/rel-style.css"
}
resource "servicenow_widget_dependency" "for_css_rel" {
  name = "tf-acc-dep-for-css-rel"
}
resource "servicenow_css_include_relation" "test" {
  css_include_id = servicenow_css_include.for_rel.id
  dependency_id  = servicenow_widget_dependency.for_css_rel.id
  order          = 100
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_css_include_relation.test"),
					resource.TestCheckResourceAttrSet("servicenow_css_include_relation.test", "css_include_id"),
					resource.TestCheckResourceAttrSet("servicenow_css_include_relation.test", "dependency_id"),
					resource.TestCheckResourceAttr("servicenow_css_include_relation.test", "order", "100"),
					resource.TestCheckResourceAttrSet("servicenow_css_include_relation.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_css_include_relation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_js_include
// ---------------------------------------------------------------------------

func TestAccResourceJsInclude_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SERVICE_PORTAL") == "" {
		t.Skip("SERVICENOW_SERVICE_PORTAL must be set to run Service Portal acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_js_include" "test" {
  display_name = "tf-acc-js-include"
  source       = "url"
  url          = "https://example.com/script.js"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_js_include.test"),
					resource.TestCheckResourceAttr("servicenow_js_include.test", "display_name", "tf-acc-js-include"),
					resource.TestCheckResourceAttr("servicenow_js_include.test", "url", "https://example.com/script.js"),
					resource.TestCheckResourceAttrSet("servicenow_js_include.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_js_include.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_js_include_relation
// ---------------------------------------------------------------------------

func TestAccResourceJsIncludeRelation_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SERVICE_PORTAL") == "" {
		t.Skip("SERVICENOW_SERVICE_PORTAL must be set to run Service Portal acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_js_include" "for_rel" {
  display_name = "tf-acc-js-for-rel"
  source       = "url"
  url          = "https://example.com/rel-script.js"
}
resource "servicenow_widget_dependency" "for_js_rel" {
  name = "tf-acc-dep-for-js-rel"
}
resource "servicenow_js_include_relation" "test" {
  js_include_id = servicenow_js_include.for_rel.id
  dependency_id = servicenow_widget_dependency.for_js_rel.id
  order         = 100
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_js_include_relation.test"),
					resource.TestCheckResourceAttrSet("servicenow_js_include_relation.test", "js_include_id"),
					resource.TestCheckResourceAttrSet("servicenow_js_include_relation.test", "dependency_id"),
					resource.TestCheckResourceAttr("servicenow_js_include_relation.test", "order", "100"),
					resource.TestCheckResourceAttrSet("servicenow_js_include_relation.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_js_include_relation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_content_css
// ---------------------------------------------------------------------------

func TestAccResourceContentCSS_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_content_css" "test" {
  name  = "tf-acc-content-css"
  type  = "local"
  style = "body { color: red; }"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_content_css.test"),
					resource.TestCheckResourceAttr("servicenow_content_css.test", "name", "tf-acc-content-css"),
					resource.TestCheckResourceAttr("servicenow_content_css.test", "type", "local"),
					resource.TestCheckResourceAttrSet("servicenow_content_css.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_content_css.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_application_menu
// ---------------------------------------------------------------------------

func TestAccResourceApplicationMenu_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_application_menu" "test" {
  title  = "TF Acc Test Menu"
  order  = 1000
  active = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_application_menu.test"),
					resource.TestCheckResourceAttr("servicenow_application_menu.test", "title", "TF Acc Test Menu"),
					resource.TestCheckResourceAttr("servicenow_application_menu.test", "order", "1000"),
					resource.TestCheckResourceAttr("servicenow_application_menu.test", "active", "true"),
					resource.TestCheckResourceAttrSet("servicenow_application_menu.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_application_menu.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_application_module
// ---------------------------------------------------------------------------

func TestAccResourceApplicationModule_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_application_menu" "for_module" {
  title  = "TF Acc Menu For Module"
  order  = 999
  active = true
}
resource "servicenow_application_module" "test" {
  title               = "TF Acc Test Module"
  application_menu_id = servicenow_application_menu.for_module.id
  order               = 100
  link_type           = "DIRECT"
  arguments           = "home.do"
  active              = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_application_module.test"),
					resource.TestCheckResourceAttr("servicenow_application_module.test", "title", "TF Acc Test Module"),
					resource.TestCheckResourceAttr("servicenow_application_module.test", "link_type", "DIRECT"),
					resource.TestCheckResourceAttrSet("servicenow_application_module.test", "application_menu_id"),
					resource.TestCheckResourceAttrSet("servicenow_application_module.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_application_module.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_system_property_relation
// ---------------------------------------------------------------------------

func TestAccResourceSystemPropertyRelation_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_system_property" "for_rel" {
  suffix = "tf.acc.prop.for.rel"
  type   = "string"
}
resource "servicenow_system_property_category" "for_rel" {
  name = "TF Acc Category For Rel"
}
resource "servicenow_system_property_relation" "test" {
  property_id = servicenow_system_property.for_rel.id
  category_id = servicenow_system_property_category.for_rel.id
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_system_property_relation.test"),
					resource.TestCheckResourceAttrSet("servicenow_system_property_relation.test", "property_id"),
					resource.TestCheckResourceAttrSet("servicenow_system_property_relation.test", "category_id"),
					resource.TestCheckResourceAttrSet("servicenow_system_property_relation.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_system_property_relation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_db_table
// Creating tables requires elevated permissions; skip gracefully when needed.
// ---------------------------------------------------------------------------

func TestAccResourceDBTable_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_DB_TABLE_TESTS") == "" {
		t.Skip("SERVICENOW_DB_TABLE_TESTS must be set to run db_table acceptance tests (requires table creation permissions)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				// ServiceNow normalizes the label to Title Case after creation
				// (e.g. "TF" -> "Tf"). Use a label that is already in Title Case
				// so we have a stable expected value.
				Config: providerBlock() + `
resource "servicenow_db_table" "test" {
  label     = "Acceptance Test Table"
  user_role = ""
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_db_table.test"),
					resource.TestCheckResourceAttr("servicenow_db_table.test", "label", "Acceptance Test Table"),
					resource.TestCheckResourceAttrSet("servicenow_db_table.test", "id"),
					resource.TestCheckResourceAttrSet("servicenow_db_table.test", "name"),
				),
			},
			{
				ResourceName:      "servicenow_db_table.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_role
// ---------------------------------------------------------------------------

func TestAccDataSourceRole_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_role" "test_ds" {
  suffix = "tf_acc_ds_test_role"
}
data "servicenow_role" "test" {
  suffix = servicenow_role.test_ds.suffix
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.servicenow_role.test", "id"),
					resource.TestCheckResourceAttr("data.servicenow_role.test", "suffix", "tf_acc_ds_test_role"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_application
// ---------------------------------------------------------------------------

func TestAccDataSourceApplication_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_application" "test_ds" {
  name    = "TF Acc DS Test Application"
  scope   = "x_tf_acc_ds_app"
  version = "1.0.0"
}
data "servicenow_application" "test" {
  name = servicenow_application.test_ds.name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.servicenow_application.test", "id"),
					resource.TestCheckResourceAttr("data.servicenow_application.test", "name", "TF Acc DS Test Application"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_application_category
// Uses a built-in ServiceNow category. Category names vary by instance.
// ---------------------------------------------------------------------------

func TestAccDataSourceApplicationCategory_basic(t *testing.T) {
	// This test requires a known application category to exist in the ServiceNow instance.
	// Developer instances may not have the expected default categories.
	// Skip if the environment variable is not explicitly set.
	if os.Getenv("SERVICENOW_APP_CATEGORY_NAME") == "" {
		t.Skip("SERVICENOW_APP_CATEGORY_NAME must be set to a known application category name to run this test")
	}

	categoryName := os.Getenv("SERVICENOW_APP_CATEGORY_NAME")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + fmt.Sprintf(`
data "servicenow_application_category" "test" {
  name = %q
}
`, categoryName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.servicenow_application_category.test", "name", categoryName),
					resource.TestCheckResourceAttrSet("data.servicenow_application_category.test", "id"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_db_table
// Uses built-in ServiceNow table "sys_user" which always exists.
// ---------------------------------------------------------------------------

func TestAccDataSourceDBTable_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
data "servicenow_db_table" "test" {
  name = "sys_user"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.servicenow_db_table.test", "name", "sys_user"),
					resource.TestCheckResourceAttrSet("data.servicenow_db_table.test", "id"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_service_catalog
// Creates a catalog, then looks it up by title.
// ---------------------------------------------------------------------------

func TestAccDataSourceServiceCatalog_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_service_catalog" "test_ds" {
  title  = "TF Acc DS Test Catalog"
  active = true
}
data "servicenow_service_catalog" "test" {
  title = servicenow_service_catalog.test_ds.title
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.servicenow_service_catalog.test", "title", "TF Acc DS Test Catalog"),
					resource.TestCheckResourceAttrSet("data.servicenow_service_catalog.test", "id"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_service_catalog_category
// Creates catalog + category, looks up category by title.
// ---------------------------------------------------------------------------

func TestAccDataSourceServiceCatalogCategory_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_service_catalog" "for_cat_ds" {
  title  = "TF Acc Catalog For Category DS"
  active = true
}
resource "servicenow_service_catalog_category" "test_ds" {
  title   = "TF Acc DS Test Category"
  catalog = servicenow_service_catalog.for_cat_ds.id
  active  = true
}
data "servicenow_service_catalog_category" "test" {
  title = servicenow_service_catalog_category.test_ds.title
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.servicenow_service_catalog_category.test", "title", "TF Acc DS Test Category"),
					resource.TestCheckResourceAttrSet("data.servicenow_service_catalog_category.test", "id"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_system_property
// Creates a system_property then looks it up by name.
// ---------------------------------------------------------------------------

func TestAccDataSourceSystemProperty_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_system_property" "test_ds" {
  suffix = "tf.acc.ds.sysprop.lookup"
  type   = "string"
}
data "servicenow_system_property" "test" {
  suffix = servicenow_system_property.test_ds.suffix
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.servicenow_system_property.test", "id"),
					resource.TestCheckResourceAttr("data.servicenow_system_property.test", "suffix", "tf.acc.ds.sysprop.lookup"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_system_property_category
// Creates a system_property_category then looks it up by name.
// ---------------------------------------------------------------------------

func TestAccDataSourceSystemPropertyCategory_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_system_property_category" "test_ds" {
  name = "TF Acc DS Test Prop Category"
}
data "servicenow_system_property_category" "test" {
  name = servicenow_system_property_category.test_ds.name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.servicenow_system_property_category.test", "name", "TF Acc DS Test Prop Category"),
					resource.TestCheckResourceAttrSet("data.servicenow_system_property_category.test", "id"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Data Source: servicenow_acl
// ServiceNow ACLs are NOT unique by name or name+operation - multiple ACL records
// can share the same name/operation for different application scopes or purposes.
// This data source is most useful when combined with a known unique ACL sys_id
// or when the user knows their instance has a uniquely-named ACL.
// Set SERVICENOW_ACL_NAME to run this test with a known-unique ACL name.
// ---------------------------------------------------------------------------

func TestAccDataSourceACL_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_ACL_NAME") == "" {
		t.Skip("SERVICENOW_ACL_NAME must be set to a uniquely-named ACL to run this test (ACLs are not globally unique by name)")
	}

	aclName := os.Getenv("SERVICENOW_ACL_NAME")
	aclOperation := os.Getenv("SERVICENOW_ACL_OPERATION") // optional

	var config string
	if aclOperation != "" {
		config = fmt.Sprintf(`
data "servicenow_acl" "test" {
  name      = %q
  operation = %q
}
`, aclName, aclOperation)
	} else {
		config = fmt.Sprintf(`
data "servicenow_acl" "test" {
  name = %q
}
`, aclName)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.servicenow_acl.test", "name", aclName),
					resource.TestCheckResourceAttrSet("data.servicenow_acl.test", "id"),
				),
			},
		},
	})
}
