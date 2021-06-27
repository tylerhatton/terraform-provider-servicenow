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
const serviceCatalogVariableListTable = "list_table"
const serviceCatalogVariableLookupTable = "lookup_table"
const serviceCatalogVariableLookupValue = "lookup_value"
const serviceCatalogVariableReference = "reference"
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
			serviceCatalogVariableListTable: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the ServiceNow table the catalog item will use to populate its list.",
			},
			serviceCatalogVariableLookupTable: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the ServiceNow table the catalog item will use to populate its lookup list.",
			},
			serviceCatalogVariableLookupValue: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the table value the catalog item will use to populate its lookup list.",
			},
			serviceCatalogVariableReference: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the table the catalog item will use to populate its reference list.",
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
	var typeString string
	switch serviceCatalogVariable.Type {
	case "33":
		typeString = "Attachment"
	case "12":
		typeString = "Break"
	case "7":
		typeString = "CheckBox"
	case "20":
		typeString = "Container End"
	case "24":
		typeString = "Container Split"
	case "19":
		typeString = "Container Start"
	case "14":
		typeString = "Custom"
	case "17":
		typeString = "Custom with Label"
	case "9":
		typeString = "Date"
	case "10":
		typeString = "Date/Time"
	case "29":
		typeString = "Duration"
	case "26":
		typeString = "Email"
	case "23":
		typeString = "HTML"
	case "28":
		typeString = "IP Address"
	case "11":
		typeString = "Label"
	case "21":
		typeString = "List Collector"
	case "22":
		typeString = "Lookup Multiple Choice"
	case "18":
		typeString = "Lookup Select Box"
	case "25":
		typeString = "Masked"
	case "2":
		typeString = "Multi Line Text"
	case "3":
		typeString = "Multiple Choice"
	case "4":
		typeString = "Numeric Scale"
	case "8":
		typeString = "Reference"
	case "31":
		typeString = "Requested For"
	case "32":
		typeString = "Rich Text Label"
	case "5":
		typeString = "Select Box"
	case "6":
		typeString = "Single Line Text"
	case "15":
		typeString = "UI Page"
	case "27":
		typeString = "URL"
	case "16":
		typeString = "Wide Single Line Text"
	case "1":
		typeString = "Yes/No"
	}

	data.SetId(serviceCatalogVariable.ID)
	data.Set(serviceCatalogVariableName, serviceCatalogVariable.Name)
	data.Set(serviceCatalogVariableQuestion, serviceCatalogVariable.Question)
	data.Set(serviceCatalogVariableTooltip, serviceCatalogVariable.Tooltip)
	data.Set(serviceCatalogVariableHelpTag, serviceCatalogVariable.HelpTag)
	data.Set(serviceCatalogVariableHelpText, serviceCatalogVariable.HelpText)
	data.Set(serviceCatalogVariableInstructions, serviceCatalogVariable.Instructions)
	data.Set(serviceCatalogVariableDefaultValue, serviceCatalogVariable.DefaultValue)
	data.Set(serviceCatalogVariableType, typeString)
	data.Set(serviceCatalogVariableCatalogItem, serviceCatalogVariable.CatalogItem)
	data.Set(serviceCatalogVariableOrder, serviceCatalogVariable.Order)
	data.Set(serviceCatalogVariableListTable, serviceCatalogVariable.ListTable)
	data.Set(serviceCatalogVariableLookupTable, serviceCatalogVariable.LookupTable)
	data.Set(serviceCatalogVariableLookupTable, serviceCatalogVariable.LookupValue)
	data.Set(serviceCatalogVariableReference, serviceCatalogVariable.Reference)
	data.Set(serviceCatalogVariableShowHelp, serviceCatalogVariable.ShowHelp)
	data.Set(serviceCatalogVariableMandatory, serviceCatalogVariable.Mandatory)
	data.Set(serviceCatalogVariableReadOnly, serviceCatalogVariable.ReadOnly)
	data.Set(serviceCatalogVariableHidden, serviceCatalogVariable.Hidden)
	data.Set(serviceCatalogVariableActive, serviceCatalogVariable.Active)
	data.Set(commonScope, serviceCatalogVariable.Scope)
}

func resourceToServiceCatalogVariable(data *schema.ResourceData) *client.ServiceCatalogVariable {
	var typeInt string
	switch data.Get(serviceCatalogVariableType).(string) {
	case "mobile":
		typeInt = "1"
	case "desktop":
		typeInt = "0"
	default:
		typeInt = "10"
	}

	serviceCatalogVariable := client.ServiceCatalogVariable{
		Name:         data.Get(serviceCatalogVariableName).(string),
		Question:     data.Get(serviceCatalogVariableQuestion).(string),
		Tooltip:      data.Get(serviceCatalogVariableTooltip).(string),
		HelpTag:      data.Get(serviceCatalogVariableHelpTag).(string),
		HelpText:     data.Get(serviceCatalogVariableHelpText).(string),
		Instructions: data.Get(serviceCatalogVariableInstructions).(string),
		DefaultValue: data.Get(serviceCatalogVariableDefaultValue).(string),
		Type:         typeInt,
		CatalogItem:  data.Get(serviceCatalogVariableCatalogItem).(string),
		Order:        data.Get(serviceCatalogVariableOrder).(string),
		ListTable:    data.Get(serviceCatalogVariableListTable).(string),
		LookupTable:  data.Get(serviceCatalogVariableLookupTable).(string),
		LookupValue:  data.Get(serviceCatalogVariableLookupValue).(string),
		Reference:    data.Get(serviceCatalogVariableReference).(string),
		ShowHelp:     data.Get(serviceCatalogVariableOrder).(bool),
		Mandatory:    data.Get(serviceCatalogVariableOrder).(bool),
		ReadOnly:     data.Get(serviceCatalogVariableReadOnly).(bool),
		Hidden:       data.Get(serviceCatalogVariableHidden).(bool),
		Active:       data.Get(serviceCatalogVariableActive).(bool),
	}
	serviceCatalogVariable.ID = data.Id()
	serviceCatalogVariable.Scope = data.Get(commonScope).(string)
	return &serviceCatalogVariable
}
