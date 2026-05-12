package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const commonProtectionPolicy = "protection_policy"
const commonScope = "scope"

func getProtectionPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Determines how application files are protected when downloaded or installed. Can be empty for no protection, 'read' for read-only protection or 'protected'.",
	}
}

func getScopeSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "global",
		ForceNew:    true,
		Description: "Associates a resource to a specific application ID in ServiceNow.",
	}
}

// setOnlyRequiredSchema Changes required parameters. For data sources, only one attribute is normally required and everything else is computed.
func setOnlyRequiredSchema(schema map[string]*schema.Schema, requiredName string) {
	for key, val := range schema {
		val.Computed = true
		val.Required = false
		val.Optional = false
		val.ForceNew = false
		val.Default = nil
		val.ValidateFunc = nil
		// DiffSuppressFunc is only meaningful when there is a config value to
		// compare against state; data source attributes are computed-only, so
		// reset it to prevent the SDK's InternalValidate from rejecting the
		// schema.
		val.DiffSuppressFunc = nil

		if key == requiredName {
			val.Computed = false
			val.Required = true
		}
	}
}
