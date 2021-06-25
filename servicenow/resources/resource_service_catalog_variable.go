package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const serviceCatalogVariableName = "name"
const serviceCatalogVariableQuestion = "question"
const serviceCatalogVariableTooltip = "tooltip"
const serviceCatalogVariableHelpTag = "help_tag"
const serviceCatalogVariableHelpText = "help_text"
const serviceCatalogVariableInstructions = "instructions"
const serviceCatalogVariableDefaultValue = "default_value"
const serviceCatalogVariableType = "type"
const serviceCatalogVariableCatalogItem = "catalog_item"
const serviceCatalogVariableOrder = "order"
const serviceCatalogVariableShowHelp = "show_help"
const serviceCatalogVariableMandatory = "mandatory"
const serviceCatalogVariableReadOnly = "read_only"
const serviceCatalogVariableHidden = "hidden"
const serviceCatalogVariableActive = "active"

// ResourceServiceCatalogVariable manages a service catalog variable in ServiceNow.
func ResourceServiceCatalogVariable() *schema.Resource {
	return &schema.Resource{
		Create: createResourceServiceCatalogVariable,
		Read:   readResourceServiceCatalogVariable,
		Update: updateResourceServiceCatalogVariable,
		Delete: deleteResourceServiceCatalogVariable,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			serviceCatalogVariableName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of variable that will be referenced in scripts throughout ServiceNow.",
			},
			serviceCatalogVariableQuestion: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the variable in catalog item.",
			},
			serviceCatalogVariableTooltip: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Tooltip hint of the variable in catalog item.",
			},
			serviceCatalogVariableHelpTag: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Help header of the variable in catalog item.",
			},
			serviceCatalogVariableHelpText: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Help text for the variable in catalog item.",
			},
			serviceCatalogVariableInstructions: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Additional instructions for the variable in catalog item.",
			},
			serviceCatalogVariableDefaultValue: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Default value for the variable in catalog item.",
			},
			serviceCatalogVariableType: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of variable to be used in catalog item.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					warns, errs = validateStringValue(val.(string), key, []string{
						"Attachment",
						"Break",
						"CheckBox",
						"Container End",
						"Container Split",
						"Container Start",
						"Custom",
						"Custom with Label",
						"Date",
						"Date/Time",
						"Duration",
						"Email",
						"HTML",
						"IP Address",
						"Label",
						"List Collector",
						"Lookup Multiple Choice",
						"Lookup Select Box",
						"Masked",
						"Multi Line Text",
						"Multiple Choice",
						"Numeric Scale",
						"Reference",
						"Requested For",
						"Rich Text Label",
						"Select Box",
						"Single Line Text",
						"UI Page",
						"URL",
						"Wide Single Line Text",
						"Yes/No",
					})
					return
				},
			},
			serviceCatalogVariableCatalogItem: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The sys id of the catalog item the variable will be assigned to",
			},
			serviceCatalogVariableOrder: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The sys id of the catalog item the variable will be assigned to",
			},
			serviceCatalogVariableShowHelp: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will enabling displaying additional help information with the variable.",
			},
			serviceCatalogVariableMandatory: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will make completion of the variable mandatory.",
			},
			serviceCatalogVariableReadOnly: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will make the variable read only.",
			},
			serviceCatalogVariableHidden: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to 'true', this property will hide the variable in the service catalog item.",
			},
			serviceCatalogVariableActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to 'true', this property will enable the variable in the service catalog item.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceServiceCatalogVariable(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogVariable := &client.ServiceCatalogVariable{}
	if err := snowClient.GetObject(client.EndpointServiceCatalogVariable, data.Id(), serviceCatalogVariable); err != nil {
		data.SetId("")
		return err
	}

	resourceFromServiceCatalogVariable(data, serviceCatalogVariable)

	return nil
}

func createResourceServiceCatalogVariable(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	serviceCatalogVariable := resourceToServiceCatalogVariable(data)
	if err := snowClient.CreateObject(client.EndpointServiceCatalogVariable, serviceCatalogVariable); err != nil {
		return err
	}

	resourceFromServiceCatalogVariable(data, serviceCatalogVariable)

	return readResourceServiceCatalogVariable(data, serviceNowClient)
}

func updateResourceServiceCatalogVariable(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointServiceCatalogVariable, resourceToServiceCatalogVariable(data)); err != nil {
		return err
	}

	return readResourceServiceCatalogVariable(data, serviceNowClient)
}

func deleteResourceServiceCatalogVariable(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return snowClient.DeleteObject(client.EndpointServiceCatalogVariable, data.Id())
}

func resourceFromServiceCatalogVariable(data *schema.ResourceData, serviceCatalogVariable *client.ServiceCatalogVariable) {
	data.SetId(serviceCatalogVariable.ID)
	data.Set(systemPropertyCategoryName, serviceCatalogVariable.Name)
	data.Set(systemPropertyCategoryTitleHTML, serviceCatalogVariable.TitleHTML)
	data.Set(commonScope, serviceCatalogVariable.Scope)
}

func resourceToServiceCatalogVariable(data *schema.ResourceData) *client.ServiceCatalogVariable {
	serviceCatalogVariable := client.ServiceCatalogVariable{
		Name:      data.Get(systemPropertyCategoryName).(string),
		TitleHTML: data.Get(systemPropertyCategoryTitleHTML).(string),
	}
	serviceCatalogVariable.ID = data.Id()
	serviceCatalogVariable.Scope = data.Get(commonScope).(string)
	return &serviceCatalogVariable
}
