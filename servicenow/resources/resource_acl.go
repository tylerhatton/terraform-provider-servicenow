package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const aclRoles = "roles"

// ResourceACL manages a security ACL (sys_security_acl) record in ServiceNow.
func ResourceACL() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_acl` manages a security access control rule (sys_security_acl) in ServiceNow.",

		CreateContext: createResourceACL,
		ReadContext:   readResourceACL,
		UpdateContext: updateResourceACL,
		DeleteContext: deleteResourceACL,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			aclName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the object being secured. Use the record/table name (e.g. 'sys_user'), or 'table.field' to secure a specific field, or 'table.*' to secure every field on the table.",
			},
			aclOperation: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Operation this ACL rule applies to. Valid values are 'read', 'write', 'create', 'delete', and 'execute'.",
				ValidateFunc: validation.StringInSlice([]string{"read", "write", "create", "delete", "execute"}, false),
			},
			aclType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "record",
				Description:  "Kind of object this ACL rule secures. Common values are 'record', 'operation', 'processor', 'client_callable_script_include', 'rest_endpoint', or 'ui_page'.",
				ValidateFunc: validation.StringInSlice([]string{"record", "operation", "processor", "client_callable_script_include", "rest_endpoint", "ui_page", "report_view_access"}, false),
			},
			aclAdminOverrides: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, users with the 'admin' role bypass this ACL rule.",
			},
			aclDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A description of the object or permissions this ACL rule secures.",
			},
			aclActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, the ACL rule is enforced.",
			},
			aclAdvanced: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, exposes the 'script' field so a custom evaluation script can be supplied.",
			},
			aclCondition: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "An encoded condition string that must evaluate to true for the rule to grant access.",
			},
			aclScript: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A server side script that must return true for the rule to grant access. Only evaluated when 'advanced' is true.",
			},
			aclRoles: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma separated list of role sys_ids whose holders are granted access by this rule.",
			},
			commonProtectionPolicy: getProtectionPolicySchema(),
			commonScope:            getScopeSchema(),
		},
	}
}

func readResourceACL(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	acl := &client.ACL{}
	if err := snowClient.GetObject(ctx, client.EndpointACL, data.Id(), acl); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromACLResource(data, acl)

	return nil
}

func createResourceACL(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	acl := resourceToACL(data)
	if err := snowClient.CreateObject(ctx, client.EndpointACL, acl); err != nil {
		return diag.FromErr(err)
	}

	resourceFromACLResource(data, acl)

	return readResourceACL(ctx, data, serviceNowClient)
}

func updateResourceACL(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointACL, resourceToACL(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceACL(ctx, data, serviceNowClient)
}

func deleteResourceACL(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointACL, data.Id()))
}

func resourceFromACLResource(data *schema.ResourceData, acl *client.ACL) {
	data.SetId(acl.ID)
	data.Set(aclName, acl.Name)
	data.Set(aclOperation, acl.Operation)
	data.Set(aclType, acl.Type)
	data.Set(aclAdminOverrides, acl.AdminOverrides)
	data.Set(aclDescription, acl.Description)
	data.Set(aclActive, acl.Active)
	data.Set(aclAdvanced, acl.Advanced)
	data.Set(aclCondition, acl.Condition)
	data.Set(aclScript, acl.Script)
	data.Set(aclRoles, acl.Roles)
	data.Set(commonProtectionPolicy, acl.ProtectionPolicy)
	data.Set(commonScope, acl.Scope)
}

func resourceToACL(data *schema.ResourceData) *client.ACL {
	acl := client.ACL{
		Name:           data.Get(aclName).(string),
		Operation:      data.Get(aclOperation).(string),
		Type:           data.Get(aclType).(string),
		AdminOverrides: data.Get(aclAdminOverrides).(bool),
		Description:    data.Get(aclDescription).(string),
		Active:         data.Get(aclActive).(bool),
		Advanced:       data.Get(aclAdvanced).(bool),
		Condition:      data.Get(aclCondition).(string),
		Script:         data.Get(aclScript).(string),
		Roles:          data.Get(aclRoles).(string),
	}
	acl.ID = data.Id()
	acl.ProtectionPolicy = data.Get(commonProtectionPolicy).(string)
	acl.Scope = data.Get(commonScope).(string)
	return &acl
}
