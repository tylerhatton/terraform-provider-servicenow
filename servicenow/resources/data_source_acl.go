package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const aclType = "type"
const aclOperation = "operation"
const aclAdminOverrides = "admin_overrides"
const aclName = "name"
const aclDescription = "description"
const aclActive = "active"
const aclAdvanced = "advanced"
const aclCondition = "condition"
const aclScript = "script"

// DataSourceACL reads the informations about a single ACL in ServiceNow.
func DataSourceACL() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		aclName: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Enter the name of the object being secured, either the record name or the table and field names.",
		},
		aclType: {
			Type:        schema.TypeString,
			Description: "Select what kind of object this ACL rule secures.",
			Computed:    true,
		},
		aclOperation: {
			Type:        schema.TypeString,
			Description: "Select the operation this ACL rule secures. Used as an additional filter when multiple ACLs share the same name (e.g. 'read', 'write', 'delete', 'create').",
			Optional:    true,
			Computed:    true,
		},
		aclAdminOverrides: {
			Type:        schema.TypeBool,
			Description: "Users with admin override this rule",
			Computed:    true,
		},
		aclDescription: {
			Type:        schema.TypeString,
			Description: "Enter a description of the object or permissions this ACL rule secures.",
			Computed:    true,
		},
		aclActive: {
			Type:        schema.TypeBool,
			Description: "Activates the ACL rule.",
			Computed:    true,
		},
		aclAdvanced: {
			Type:        schema.TypeBool,
			Description: "Displays the Script field when active.",
			Computed:    true,
		},
		aclCondition: {
			Type:        schema.TypeString,
			Description: "Selects the fields and values that must be true for users to access the object.",
			Computed:    true,
		},
		aclScript: {
			Type:        schema.TypeString,
			Description: "Custom script describing the permissions required to access the object.",
			Computed:    true,
		},
		commonProtectionPolicy: getProtectionPolicySchema(),
		commonScope:            getScopeSchema(),
	}

	setOnlyRequiredSchema(resourceSchema, aclName)

	// After setOnlyRequiredSchema forces everything to Computed, restore operation and type
	// as Optional+Computed so users can provide them to disambiguate ACLs with the same name.
	resourceSchema[aclOperation].Optional = true
	resourceSchema[aclType].Optional = true

	return &schema.Resource{
		Description: "`servicenow_acl` data source can be used to retrieve information of a single ACL in ServiceNow. Use operation and type to disambiguate when multiple ACLs share the same name.",
		ReadContext: readDataSourceACL,
		Schema:      resourceSchema,
	}
}

func readDataSourceACL(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	acl := &client.ACL{}
	name := data.Get(aclName).(string)
	operation := data.Get(aclOperation).(string)
	aclTypeVal := data.Get(aclType).(string)

	var err error
	if operation != "" || aclTypeVal != "" {
		// Build compound query to uniquely identify the ACL.
		query := "name=" + name
		if operation != "" {
			query += "^operation=" + operation
		}
		if aclTypeVal != "" {
			query += "^type=" + aclTypeVal
		}
		err = snowClient.GetObjectByQuery(client.EndpointACL, query, acl)
	} else {
		err = snowClient.GetObjectByName(client.EndpointACL, name, acl)
	}
	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromACL(data, acl)

	return nil
}

func resourceFromACL(data *schema.ResourceData, acl *client.ACL) {
	data.SetId(acl.ID)
	data.Set(aclType, acl.Type)
	data.Set(aclOperation, acl.Operation)
	data.Set(aclAdminOverrides, acl.AdminOverrides)
	data.Set(aclName, acl.Name)
	data.Set(aclDescription, acl.Description)
	data.Set(aclActive, acl.Active)
	data.Set(aclAdvanced, acl.Advanced)
	data.Set(aclCondition, acl.Condition)
	data.Set(aclScript, acl.Script)
	data.Set(commonProtectionPolicy, acl.ProtectionPolicy)
	data.Set(commonScope, acl.Scope)
}
