package resources_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ---------------------------------------------------------------------------
// Resource: servicenow_user
// ---------------------------------------------------------------------------

func TestAccUser_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_user" "test" {
  user_name  = "tf-acc-user-basic"
  first_name = "Test"
  last_name  = "User"
  email      = "tf-acc-user@example.com"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_user.test"),
					resource.TestCheckResourceAttr("servicenow_user.test", "user_name", "tf-acc-user-basic"),
					resource.TestCheckResourceAttr("servicenow_user.test", "first_name", "Test"),
					resource.TestCheckResourceAttr("servicenow_user.test", "last_name", "User"),
					resource.TestCheckResourceAttrSet("servicenow_user.test", "id"),
				),
			},
			{
				ResourceName:            "servicenow_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_password"},
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_group
// ---------------------------------------------------------------------------

func TestAccGroup_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_group" "test" {
  name        = "tf-acc-group-basic"
  description = "TF acceptance test group"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_group.test"),
					resource.TestCheckResourceAttr("servicenow_group.test", "name", "tf-acc-group-basic"),
					resource.TestCheckResourceAttr("servicenow_group.test", "description", "TF acceptance test group"),
					resource.TestCheckResourceAttrSet("servicenow_group.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_user_role
// NOTE: sys_user_has_role is gated by the security_admin elevated role. The
// JSONv2 API call as admin returns "Insufficient rights to insert a
// sys_user_has_role record" because security_admin must be elevated through
// the UI. Skip by default; opt-in via SERVICENOW_SECURITY_ADMIN=1 on instances
// where the credential has the elevated role.
// ---------------------------------------------------------------------------

func TestAccUserRole_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SECURITY_ADMIN") == "" {
		t.Skip("SERVICENOW_SECURITY_ADMIN must be set to run user_role acceptance tests (sys_user_has_role inserts require the security_admin role to be elevated)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_user" "for_user_role" {
  user_name = "tf-acc-user-for-role"
  last_name = "RoleUser"
}
resource "servicenow_role" "for_user_role" {
  suffix = "tf_acc_user_role_test"
}
resource "servicenow_user_role" "test" {
  user = servicenow_user.for_user_role.id
  role = servicenow_role.for_user_role.id
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_user_role.test"),
					resource.TestCheckResourceAttrSet("servicenow_user_role.test", "user"),
					resource.TestCheckResourceAttrSet("servicenow_user_role.test", "role"),
					resource.TestCheckResourceAttrSet("servicenow_user_role.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_user_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_group_member
// ---------------------------------------------------------------------------

func TestAccGroupMember_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_user" "for_member" {
  user_name = "tf-acc-user-for-member"
  last_name = "MemberUser"
}
resource "servicenow_group" "for_member" {
  name = "tf-acc-group-for-member"
}
resource "servicenow_group_member" "test" {
  user  = servicenow_user.for_member.id
  group = servicenow_group.for_member.id
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_group_member.test"),
					resource.TestCheckResourceAttrSet("servicenow_group_member.test", "user"),
					resource.TestCheckResourceAttrSet("servicenow_group_member.test", "group"),
					resource.TestCheckResourceAttrSet("servicenow_group_member.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_group_member.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_group_role
// ---------------------------------------------------------------------------

func TestAccGroupRole_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_group" "for_grouprole" {
  name = "tf-acc-group-for-role"
}
resource "servicenow_role" "for_grouprole" {
  suffix = "tf_acc_group_role_test"
}
resource "servicenow_group_role" "test" {
  group = servicenow_group.for_grouprole.id
  role  = servicenow_role.for_grouprole.id
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_group_role.test"),
					resource.TestCheckResourceAttrSet("servicenow_group_role.test", "group"),
					resource.TestCheckResourceAttrSet("servicenow_group_role.test", "role"),
					resource.TestCheckResourceAttrSet("servicenow_group_role.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_group_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_business_rule
// ---------------------------------------------------------------------------

func TestAccBusinessRule_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_business_rule" "test" {
  name          = "tf-acc-business-rule"
  table         = "incident"
  when          = "before"
  action_insert = true
  script        = "(function executeRule(current, previous) { /* test */ })(current, previous);"
  description   = "TF acceptance test business rule"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_business_rule.test"),
					resource.TestCheckResourceAttr("servicenow_business_rule.test", "name", "tf-acc-business-rule"),
					resource.TestCheckResourceAttr("servicenow_business_rule.test", "table", "incident"),
					resource.TestCheckResourceAttr("servicenow_business_rule.test", "when", "before"),
					resource.TestCheckResourceAttrSet("servicenow_business_rule.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_business_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_client_script
// ---------------------------------------------------------------------------

func TestAccClientScript_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_client_script" "test" {
  name        = "tf-acc-client-script"
  table       = "incident"
  type        = "onLoad"
  script      = "function onLoad() { /* test */ }"
  description = "TF acceptance test client script"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_client_script.test"),
					resource.TestCheckResourceAttr("servicenow_client_script.test", "name", "tf-acc-client-script"),
					resource.TestCheckResourceAttr("servicenow_client_script.test", "table", "incident"),
					resource.TestCheckResourceAttr("servicenow_client_script.test", "type", "onLoad"),
					resource.TestCheckResourceAttrSet("servicenow_client_script.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_client_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_ui_action
// ---------------------------------------------------------------------------

func TestAccUIAction_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_ui_action" "test" {
  name        = "tf-acc-ui-action"
  table       = "incident"
  action_name = "tf_acc_ui_action_test"
  form_button = true
  script      = "// test UI action"
  comments    = "TF acceptance test UI action"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_ui_action.test"),
					resource.TestCheckResourceAttr("servicenow_ui_action.test", "name", "tf-acc-ui-action"),
					resource.TestCheckResourceAttr("servicenow_ui_action.test", "table", "incident"),
					resource.TestCheckResourceAttrSet("servicenow_ui_action.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_ui_action.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_ui_policy
// ---------------------------------------------------------------------------

func TestAccUIPolicy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_ui_policy" "test" {
  short_description = "tf-acc-ui-policy"
  table             = "incident"
  description       = "TF acceptance test UI policy"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_ui_policy.test"),
					resource.TestCheckResourceAttr("servicenow_ui_policy.test", "short_description", "tf-acc-ui-policy"),
					resource.TestCheckResourceAttr("servicenow_ui_policy.test", "table", "incident"),
					resource.TestCheckResourceAttrSet("servicenow_ui_policy.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_ui_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_ui_policy_action
// NOTE: Some ServiceNow instances ship with HasNobodyRole ACLs on
// sys_ui_policy_action.ui_policy that prevent the parent UI policy reference
// from being set via the JSONv2 API. dev391819 has these in place. Skip the
// test by default and let installations that have lifted those ACLs run it
// explicitly with SERVICENOW_UI_POLICY_ACTION_TESTS=1.
// ---------------------------------------------------------------------------

func TestAccUIPolicyAction_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_UI_POLICY_ACTION_TESTS") == "" {
		t.Skip("SERVICENOW_UI_POLICY_ACTION_TESTS must be set to run ui_policy_action acceptance tests (the ui_policy reference is blocked by default ACLs on stock instances)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_ui_policy" "parent" {
  short_description = "tf-acc-ui-policy-for-action"
  table             = "incident"
}
resource "servicenow_ui_policy_action" "test" {
  ui_policy  = servicenow_ui_policy.parent.id
  field_name = "short_description"
  mandatory  = "true"
  visible    = "true"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_ui_policy_action.test"),
					resource.TestCheckResourceAttrSet("servicenow_ui_policy_action.test", "ui_policy"),
					resource.TestCheckResourceAttr("servicenow_ui_policy_action.test", "field_name", "short_description"),
					resource.TestCheckResourceAttrSet("servicenow_ui_policy_action.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_ui_policy_action.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_dictionary
// Adds a custom column to the incident table. Uses u_ prefix per SN convention.
// ---------------------------------------------------------------------------

func TestAccDictionary_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_dictionary" "test" {
  name          = "incident"
  element       = "u_tf_acc_dict_field"
  column_label  = "TF Acc Dict Field"
  internal_type = "string"
  max_length    = 40
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_dictionary.test"),
					resource.TestCheckResourceAttr("servicenow_dictionary.test", "name", "incident"),
					resource.TestCheckResourceAttr("servicenow_dictionary.test", "element", "u_tf_acc_dict_field"),
					resource.TestCheckResourceAttr("servicenow_dictionary.test", "column_label", "TF Acc Dict Field"),
					resource.TestCheckResourceAttrSet("servicenow_dictionary.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_dictionary.test",
				ImportState:       true,
				ImportStateVerify: true,
				// ServiceNow returns internal_type as a sys_glide_object sys_id for
				// freshly imported records. The DiffSuppressFunc on the resource
				// hides the planning-time delta, but import-time state comparison
				// still sees the sys_id form versus the user-supplied "string".
				ImportStateVerifyIgnore: []string{"internal_type"},
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_choice
// Adds a choice value to an existing column on incident.
// ---------------------------------------------------------------------------

func TestAccChoice_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_dictionary" "for_choice" {
  name          = "incident"
  element       = "u_tf_acc_choice_field"
  column_label  = "TF Acc Choice Field"
  internal_type = "string"
  max_length    = 40
  choice        = 3
}
resource "servicenow_choice" "test" {
  name     = "incident"
  element  = servicenow_dictionary.for_choice.element
  value    = "tf_acc_choice_val"
  label    = "TF Acc Choice"
  sequence = 100
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_choice.test"),
					resource.TestCheckResourceAttr("servicenow_choice.test", "name", "incident"),
					resource.TestCheckResourceAttr("servicenow_choice.test", "element", "u_tf_acc_choice_field"),
					resource.TestCheckResourceAttr("servicenow_choice.test", "value", "tf_acc_choice_val"),
					resource.TestCheckResourceAttrSet("servicenow_choice.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_choice.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_acl
// NOTE: sys_security_acl is gated by the security_admin elevated role. The
// JSONv2 API call as admin returns "Insufficient rights to insert a
// sys_security_acl record" because security_admin must be elevated through
// the UI. Skip by default; opt-in via SERVICENOW_SECURITY_ADMIN=1 on
// instances where the credential has the elevated role.
// ---------------------------------------------------------------------------

func TestAccACL_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_SECURITY_ADMIN") == "" {
		t.Skip("SERVICENOW_SECURITY_ADMIN must be set to run servicenow_acl acceptance tests (sys_security_acl inserts require the security_admin role to be elevated)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_acl" "test" {
  name        = "tf_acc_acl_test_resource"
  operation   = "read"
  type        = "record"
  description = "TF acceptance test ACL"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_acl.test"),
					resource.TestCheckResourceAttr("servicenow_acl.test", "name", "tf_acc_acl_test_resource"),
					resource.TestCheckResourceAttr("servicenow_acl.test", "operation", "read"),
					resource.TestCheckResourceAttr("servicenow_acl.test", "type", "record"),
					resource.TestCheckResourceAttrSet("servicenow_acl.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_acl.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_notification
// ---------------------------------------------------------------------------

func TestAccNotification_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_notification" "test" {
  name         = "tf-acc-notification"
  table        = "incident"
  subject      = "TF Acc Test Notification"
  message_html = "<p>Test message body</p>"
  description  = "TF acceptance test notification"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_notification.test"),
					resource.TestCheckResourceAttr("servicenow_notification.test", "name", "tf-acc-notification"),
					resource.TestCheckResourceAttr("servicenow_notification.test", "table", "incident"),
					resource.TestCheckResourceAttrSet("servicenow_notification.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_notification.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_email_template
// ---------------------------------------------------------------------------

func TestAccEmailTemplate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_email_template" "test" {
  name         = "tf-acc-email-template"
  table        = "incident"
  subject      = "TF Acc Email Template Subject"
  message_html = "<p>TF acceptance test email template body</p>"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_email_template.test"),
					resource.TestCheckResourceAttr("servicenow_email_template.test", "name", "tf-acc-email-template"),
					resource.TestCheckResourceAttr("servicenow_email_template.test", "subject", "TF Acc Email Template Subject"),
					resource.TestCheckResourceAttrSet("servicenow_email_template.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_email_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_scheduled_job
// ---------------------------------------------------------------------------

func TestAccScheduledJob_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_scheduled_job" "test" {
  name     = "tf-acc-scheduled-job"
  script   = "gs.info('tf-acc scheduled job');"
  run_type = "daily"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_scheduled_job.test"),
					resource.TestCheckResourceAttr("servicenow_scheduled_job.test", "name", "tf-acc-scheduled-job"),
					resource.TestCheckResourceAttr("servicenow_scheduled_job.test", "run_type", "daily"),
					resource.TestCheckResourceAttrSet("servicenow_scheduled_job.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_scheduled_job.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_script_action
// ---------------------------------------------------------------------------

func TestAccScriptAction_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_script_action" "test" {
  name        = "tf-acc-script-action"
  event_name  = "tf.acc.event"
  script      = "gs.info('tf-acc script action');"
  description = "TF acceptance test script action"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_script_action.test"),
					resource.TestCheckResourceAttr("servicenow_script_action.test", "name", "tf-acc-script-action"),
					resource.TestCheckResourceAttr("servicenow_script_action.test", "event_name", "tf.acc.event"),
					resource.TestCheckResourceAttrSet("servicenow_script_action.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_script_action.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_assignment_rule
// ---------------------------------------------------------------------------

func TestAccAssignmentRule_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_assignment_rule" "test" {
  name        = "tf-acc-assignment-rule"
  table       = "incident"
  description = "TF acceptance test assignment rule"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_assignment_rule.test"),
					resource.TestCheckResourceAttr("servicenow_assignment_rule.test", "name", "tf-acc-assignment-rule"),
					resource.TestCheckResourceAttr("servicenow_assignment_rule.test", "table", "incident"),
					resource.TestCheckResourceAttrSet("servicenow_assignment_rule.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_assignment_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_data_lookup
// ---------------------------------------------------------------------------

func TestAccDataLookup_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_data_lookup" "test" {
  name         = "tf-acc-data-lookup"
  table        = "incident"
  lookup_table = "dl_matcher"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_data_lookup.test"),
					resource.TestCheckResourceAttr("servicenow_data_lookup.test", "name", "tf-acc-data-lookup"),
					resource.TestCheckResourceAttr("servicenow_data_lookup.test", "table", "incident"),
					resource.TestCheckResourceAttrSet("servicenow_data_lookup.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_data_lookup.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_transform_map
// ---------------------------------------------------------------------------

func TestAccTransformMap_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				// ServiceNow ships a "Transform Validator" business rule that
				// rejects transform maps where source_table == target_table, so
				// we use two distinct tables here.
				Config: providerBlock() + `
resource "servicenow_transform_map" "test" {
  name         = "tf-acc-transform-map"
  source_table = "sys_user"
  target_table = "sys_user_group"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_transform_map.test"),
					resource.TestCheckResourceAttr("servicenow_transform_map.test", "name", "tf-acc-transform-map"),
					resource.TestCheckResourceAttr("servicenow_transform_map.test", "source_table", "sys_user"),
					resource.TestCheckResourceAttr("servicenow_transform_map.test", "target_table", "sys_user_group"),
					resource.TestCheckResourceAttrSet("servicenow_transform_map.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_transform_map.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_transform_entry
// ---------------------------------------------------------------------------

func TestAccTransformEntry_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Source and target tables must differ to satisfy the
				// Transform Validator business rule. Use sys_user -> sys_user_group
				// and map "name" on both sides (a column that exists on the
				// target table).
				Config: providerBlock() + `
resource "servicenow_transform_map" "for_entry" {
  name         = "tf-acc-tm-for-entry"
  source_table = "sys_user"
  target_table = "sys_user_group"
}
resource "servicenow_transform_entry" "test" {
  map          = servicenow_transform_map.for_entry.id
  source_field = "user_name"
  target_field = "name"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_transform_entry.test"),
					resource.TestCheckResourceAttrSet("servicenow_transform_entry.test", "map"),
					resource.TestCheckResourceAttr("servicenow_transform_entry.test", "source_field", "user_name"),
					resource.TestCheckResourceAttr("servicenow_transform_entry.test", "target_field", "name"),
					resource.TestCheckResourceAttrSet("servicenow_transform_entry.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_transform_entry.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_certificate
// Uses a real self-signed test PEM certificate.
// ---------------------------------------------------------------------------

// A self-signed test certificate (PEM) used for acceptance tests. Generated
// from a throwaway test key; expires far in the future.
const tfAccTestPEM = `-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUYjwR1ZQ8RnIxNL+ojSydsfQHy+0wDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCVVMxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMzAxMDEwMDAwMDBaFw0zMzAx
MDEwMDAwMDBaMEUxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDC5SJgr1xa+r3KFm1J9Eb37Khy7zEhBgrR8a1S/PdL
2gxL9j0AOcjnyaLI5T2qZJZ3uSbU0EorL8wDxhPjyqcdHN1ml8h7Vmqlzc6Wxo9Q
RPxXh38i5dG7e0Z3w5cQyzhRz1RoTjV6QrPIzWiK+TGOIm+gIm6OBkV0eEZw8GfP
P1jcXJpoJgUMRQTOByh4PaZ0/MV2yMaJ5RVfYZWlBgYxTNoQXBdjg2Eh4QlQyhqz
J7CSWqlnh3OtVfvBgg4SXt5Yqd5RGdRcGGGZJ4o3IhYRQjJ8I3I1NlSRqZjV/Lcs
JU8e8wfbeLk0AT9TwLg+0nDvR/SUNUSk2y2eVk24AfZ/AgMBAAGjUzBRMB0GA1Ud
DgQWBBSqK0nXJDQjJC3LM3oVwQzZxRzGsTAfBgNVHSMEGDAWgBSqK0nXJDQjJC3L
M3oVwQzZxRzGsTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAA
abkbjMaIzGcKMzMOLqzS5xPYxiF7twvyAGFiZ9XSdqWXfdY3R3O2NRwsCgwwBgYW
EuLkX0J4qY8gWdczEOMTvUCNDqu7DgIxZNcD5L4DcVHwsCgKZBgZchsxsYY+JePc
KO6mh3Tl5d6gBp9Wg2QLkBgcJaQltGcShVqYAYsBOL3oFDU/8YQKKjMcEr5LZ2bH
+kfBQA8e2OCKvJh6JeOEYqJBnp4zKvJ5d1S1JEbRdR0dF5l5kBwG7AVKwbe3w/AY
7g4ojb6q4ARtIjOq/T2/PFC2YYK7r4ZPwRJDpEcWXyZNUaZjKBFKxr0VBjZNTd9C
Ki29VgaUbjMUW1YqcKga
-----END CERTIFICATE-----`

func TestAccCertificate_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_CERT_TESTS") == "" {
		t.Skip("SERVICENOW_CERT_TESTS must be set to run servicenow_certificate acceptance tests (requires valid PEM certificate ServiceNow can parse)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_certificate" "test" {
  name              = "tf-acc-certificate"
  short_description = "TF acceptance test certificate"
  format            = "pem"
  type              = "trust_store"
  pem_certificate   = <<-EOT
` + tfAccTestPEM + `
  EOT
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_certificate.test"),
					resource.TestCheckResourceAttr("servicenow_certificate.test", "name", "tf-acc-certificate"),
					resource.TestCheckResourceAttr("servicenow_certificate.test", "format", "pem"),
					resource.TestCheckResourceAttrSet("servicenow_certificate.test", "id"),
				),
			},
			{
				ResourceName:            "servicenow_certificate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pem_certificate", "key_store_password", "key_store"},
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_mid_server
// MID server records are normally created by installing the MID server agent
// (which auto-registers). Skip this test unless explicitly enabled.
// ---------------------------------------------------------------------------

func TestAccMidServer_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_MID_SERVER_TESTS") == "" {
		t.Skip("SERVICENOW_MID_SERVER_TESTS must be set to run mid_server acceptance tests (MID servers are normally registered by the MID agent, not via API)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_mid_server" "test" {
  name      = "tf-acc-mid-server"
  host_name = "tf-acc-host"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_mid_server.test"),
					resource.TestCheckResourceAttr("servicenow_mid_server.test", "name", "tf-acc-mid-server"),
					resource.TestCheckResourceAttrSet("servicenow_mid_server.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_mid_server.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_jdbc_connection
// Requires a parent connection alias of type jdbc_connection.
// ---------------------------------------------------------------------------

func TestAccJdbcConnection_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_alias" "for_jdbc" {
  name            = "tf-acc-jdbc-alias"
  type            = "connection"
  connection_type = "jdbc_connection"
}
resource "servicenow_jdbc_connection" "test" {
  name             = "tf-acc-jdbc-connection"
  connection_alias = servicenow_alias.for_jdbc.id
  connection_url   = "jdbc:mysql://localhost:3306/tfacc"
  database_name    = "tfacc"
  database_type    = "mysql"
  active           = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_jdbc_connection.test"),
					resource.TestCheckResourceAttr("servicenow_jdbc_connection.test", "name", "tf-acc-jdbc-connection"),
					resource.TestCheckResourceAttr("servicenow_jdbc_connection.test", "connection_url", "jdbc:mysql://localhost:3306/tfacc"),
					resource.TestCheckResourceAttrSet("servicenow_jdbc_connection.test", "connection_alias"),
					resource.TestCheckResourceAttrSet("servicenow_jdbc_connection.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_jdbc_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_flow
// Flow Designer flows are normally authored in the UI. Skip unless enabled.
// ---------------------------------------------------------------------------

func TestAccFlow_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_FLOW_TESTS") == "" {
		t.Skip("SERVICENOW_FLOW_TESTS must be set to run servicenow_flow acceptance tests (Flow Designer flows are normally authored in the UI)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_flow" "test" {
  name        = "tf-acc-flow"
  description = "TF acceptance test flow"
  active      = false
  category    = "flow"
  status      = "draft"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_flow.test"),
					resource.TestCheckResourceAttr("servicenow_flow.test", "name", "tf-acc-flow"),
					resource.TestCheckResourceAttrSet("servicenow_flow.test", "id"),
				),
			},
			{
				ResourceName:      "servicenow_flow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Resource: servicenow_encryption_context
// Requires Edge Encryption plugin. Skip unless explicitly enabled.
// ---------------------------------------------------------------------------

func TestAccEncryptionContext_basic(t *testing.T) {
	if os.Getenv("SERVICENOW_EDGE_ENCRYPTION") == "" {
		t.Skip("SERVICENOW_EDGE_ENCRYPTION must be set to run encryption_context acceptance tests (requires the Edge Encryption plugin)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock() + `
resource "servicenow_encryption_context" "test" {
  name        = "tf-acc-encryption-context"
  type        = "standard"
  algorithm   = "AES_256"
  description = "TF acceptance test encryption context"
}
`,
				Check: resource.ComposeTestCheckFunc(
					checkExists("servicenow_encryption_context.test"),
					resource.TestCheckResourceAttr("servicenow_encryption_context.test", "name", "tf-acc-encryption-context"),
					resource.TestCheckResourceAttrSet("servicenow_encryption_context.test", "id"),
				),
			},
			{
				ResourceName:            "servicenow_encryption_context.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption_key"},
			},
		},
	})
}
