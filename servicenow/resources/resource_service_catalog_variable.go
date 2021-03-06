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
const serviceCatalogVariableReferenceQualifier = "reference_qualifier"
const serviceCatalogVariableShowHelp = "show_help"
const serviceCatalogVariableMandatory = "mandatory"
const serviceCatalogVariableReadOnly = "read_only"
const serviceCatalogVariableHidden = "hidden"
const serviceCatalogVariableActive = "active"

// ResourceServiceCatalogVariable manages a service catalog variable in ServiceNow.
func ResourceServiceCatalogVariable() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_service_catalog_variable` manages a service catalog variable configuration within ServiceNow.",

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
				Optional:    true,
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
				Default:     "More information",
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
			serviceCatalogVariableReferenceQualifier: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The table filter applied to the reference lookup.",
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
	data.Set(serviceCatalogVariableReferenceQualifier, serviceCatalogVariable.ReferenceQualifier)
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
	case "Attachment":
		typeInt = "33"
	case "Break":
		typeInt = "12"
	case "CheckBox":
		typeInt = "7"
	case "Container End":
		typeInt = "20"
	case "Container Split":
		typeInt = "24"
	case "Container Start":
		typeInt = "19"
	case "Custom":
		typeInt = "14"
	case "Custom with Label":
		typeInt = "17"
	case "Date":
		typeInt = "9"
	case "Date/Time":
		typeInt = "10"
	case "Duration":
		typeInt = "29"
	case "Email":
		typeInt = "26"
	case "HTML":
		typeInt = "23"
	case "IP Address":
		typeInt = "28"
	case "Label":
		typeInt = "11"
	case "List Collector":
		typeInt = "21"
	case "Lookup Multiple Choice":
		typeInt = "22"
	case "Lookup Select Box":
		typeInt = "18"
	case "Masked":
		typeInt = "25"
	case "Multi Line Text":
		typeInt = "2"
	case "Multiple Choice":
		typeInt = "3"
	case "Numeric Scale":
		typeInt = "4"
	case "Reference":
		typeInt = "8"
	case "Requested For":
		typeInt = "31"
	case "Rich Text Label":
		typeInt = "32"
	case "Select Box":
		typeInt = "5"
	case "Single Line Text":
		typeInt = "6"
	case "UI Page":
		typeInt = "15"
	case "URL":
		typeInt = "27"
	case "Wide Single Line Text":
		typeInt = "16"
	case "Yes/No":
		typeInt = "1"
	}

	serviceCatalogVariable := client.ServiceCatalogVariable{
		Name:               data.Get(serviceCatalogVariableName).(string),
		Question:           data.Get(serviceCatalogVariableQuestion).(string),
		Tooltip:            data.Get(serviceCatalogVariableTooltip).(string),
		HelpTag:            data.Get(serviceCatalogVariableHelpTag).(string),
		HelpText:           data.Get(serviceCatalogVariableHelpText).(string),
		Instructions:       data.Get(serviceCatalogVariableInstructions).(string),
		DefaultValue:       data.Get(serviceCatalogVariableDefaultValue).(string),
		Type:               typeInt,
		CatalogItem:        data.Get(serviceCatalogVariableCatalogItem).(string),
		Order:              data.Get(serviceCatalogVariableOrder).(string),
		ListTable:          data.Get(serviceCatalogVariableListTable).(string),
		LookupTable:        data.Get(serviceCatalogVariableLookupTable).(string),
		LookupValue:        data.Get(serviceCatalogVariableLookupValue).(string),
		Reference:          data.Get(serviceCatalogVariableReference).(string),
		ReferenceQualifier: data.Get(serviceCatalogVariableReferenceQualifier).(string),
		ShowHelp:           data.Get(serviceCatalogVariableShowHelp).(bool),
		Mandatory:          data.Get(serviceCatalogVariableMandatory).(bool),
		ReadOnly:           data.Get(serviceCatalogVariableReadOnly).(bool),
		Hidden:             data.Get(serviceCatalogVariableHidden).(bool),
		Active:             data.Get(serviceCatalogVariableActive).(bool),
	}
	serviceCatalogVariable.ID = data.Id()
	serviceCatalogVariable.Scope = data.Get(commonScope).(string)
	return &serviceCatalogVariable
}
