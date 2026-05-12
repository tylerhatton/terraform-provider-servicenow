package resources

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

// sysIDPattern matches a 32-character hex ServiceNow sys_id.
var sysIDPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)

// dictionaryInternalTypeSuppressDiff suppresses diffs when one side is the
// configured internal_type name (e.g. "string") and the other is the sys_id
// of the corresponding sys_glide_object record. ServiceNow's JSONv2 API
// returns the sys_id form of the reference for newly written records, even
// though the input was the type name.
func dictionaryInternalTypeSuppressDiff(_, old, new string, _ *schema.ResourceData) bool {
	if old == new {
		return true
	}
	// Treat one side being a sys_id and the other being any non-empty value as
	// equivalent. The user-configured value is the source of truth for state.
	if sysIDPattern.MatchString(old) && new != "" {
		return true
	}
	if sysIDPattern.MatchString(new) && old != "" {
		return true
	}
	return false
}

const dictionaryName = "name"
const dictionaryElement = "element"
const dictionaryColumnLabel = "column_label"
const dictionaryInternalType = "internal_type"
const dictionaryMaxLength = "max_length"
const dictionaryMandatory = "mandatory"
const dictionaryReadOnly = "read_only"
const dictionaryActive = "active"
const dictionaryDisplay = "display"
const dictionaryUnique = "unique"
const dictionaryDefaultValue = "default_value"
const dictionaryComments = "comments"
const dictionaryReference = "reference"
const dictionaryDynamicCreation = "dynamic_creation"
const dictionaryDependent = "dependent"
const dictionaryDependentOnField = "dependent_on_field"
const dictionaryChoice = "choice"

// ResourceDictionary manages a column definition (dictionary entry) in ServiceNow.
func ResourceDictionary() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_dictionary` manages a column definition (sys_dictionary entry) for a table within ServiceNow.",

		CreateContext: createResourceDictionary,
		ReadContext:   readResourceDictionary,
		UpdateContext: updateResourceDictionary,
		DeleteContext: deleteResourceDictionary,

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
			dictionaryName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of the table this column belongs to.",
			},
			dictionaryElement: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of the column (field) being defined.",
			},
			dictionaryColumnLabel: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Human readable label for the column.",
			},
			dictionaryInternalType: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "string",
				Description:      "Internal data type of the column (e.g. 'string', 'integer', 'boolean', 'glide_date', 'glide_date_time', 'reference', 'choice').",
				DiffSuppressFunc: dictionaryInternalTypeSuppressDiff,
				ValidateFunc: validation.StringInSlice([]string{
					"string",
					"integer",
					"boolean",
					"glide_date",
					"glide_date_time",
					"reference",
					"choice",
					"decimal",
					"float",
					"longint",
					"long",
					"password",
					"password2",
					"url",
					"email",
					"html",
					"script",
					"script_plain",
					"translated_text",
					"translated_html",
					"journal",
					"journal_input",
					"glide_duration",
					"currency",
					"phone_number",
					"ip_addr",
					"sys_class_name",
					"user_image",
					"image",
					"document_id",
				}, false),
			},
			dictionaryMaxLength: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The maximum length permitted for values stored in the column.",
			},
			dictionaryMandatory: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the column must contain a value.",
			},
			dictionaryReadOnly: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the column cannot be edited via the standard UI.",
			},
			dictionaryActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag indicating if this dictionary entry is active.",
			},
			dictionaryDisplay: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this column is the table's display column.",
			},
			dictionaryUnique: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, values stored in the column must be unique across rows.",
			},
			dictionaryDefaultValue: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Default value applied to the column when a new record is created.",
			},
			dictionaryComments: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Free form comments describing this column.",
			},
			dictionaryReference: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Target table name when 'internal_type' is 'reference'.",
			},
			dictionaryDynamicCreation: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this column is dynamically created at runtime.",
			},
			dictionaryDependent: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of a column on which this column's choices depend.",
			},
			dictionaryDependentOnField: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of the field this column depends on for dynamic display.",
			},
			dictionaryChoice: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "Choice type for the column. 0 = none, 1 = suggestion, 3 = dropdown.",
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3}),
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceDictionary(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	dictionary := &client.Dictionary{}
	if err := snowClient.GetObject(ctx, client.EndpointDictionary, data.Id(), dictionary); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromDictionary(data, dictionary)

	return nil
}

func createResourceDictionary(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	dictionary := resourceToDictionary(data)
	if err := snowClient.CreateObject(ctx, client.EndpointDictionary, dictionary); err != nil {
		return diag.FromErr(err)
	}

	resourceFromDictionary(data, dictionary)

	return readResourceDictionary(ctx, data, serviceNowClient)
}

func updateResourceDictionary(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointDictionary, resourceToDictionary(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceDictionary(ctx, data, serviceNowClient)
}

func deleteResourceDictionary(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointDictionary, data.Id()))
}

func resourceFromDictionary(data *schema.ResourceData, dictionary *client.Dictionary) {
	data.SetId(dictionary.ID)
	data.Set(dictionaryName, dictionary.Name)
	data.Set(dictionaryElement, dictionary.Element)
	data.Set(dictionaryColumnLabel, dictionary.ColumnLabel)
	// ServiceNow may return internal_type as a sys_id reference to
	// sys_glide_object; preserve the user-supplied value in that case so state
	// stays comparable with the configuration.
	if !sysIDPattern.MatchString(dictionary.InternalType) {
		data.Set(dictionaryInternalType, dictionary.InternalType)
	}
	maxLen, err := strconv.Atoi(dictionary.MaxLength)
	if err != nil {
		maxLen = 0
	}
	data.Set(dictionaryMaxLength, maxLen)
	data.Set(dictionaryMandatory, dictionary.Mandatory)
	data.Set(dictionaryReadOnly, dictionary.ReadOnly)
	data.Set(dictionaryActive, dictionary.Active)
	data.Set(dictionaryDisplay, dictionary.Display)
	data.Set(dictionaryUnique, dictionary.Unique)
	data.Set(dictionaryDefaultValue, dictionary.DefaultValue)
	data.Set(dictionaryComments, dictionary.Comments)
	data.Set(dictionaryReference, dictionary.Reference)
	data.Set(dictionaryDynamicCreation, dictionary.DynamicCreation)
	data.Set(dictionaryDependent, dictionary.Dependent)
	data.Set(dictionaryDependentOnField, dictionary.DependentOnField)
	choice, err := strconv.Atoi(dictionary.Choice)
	if err != nil {
		choice = 0
	}
	data.Set(dictionaryChoice, choice)
	data.Set(commonScope, dictionary.Scope)
}

func resourceToDictionary(data *schema.ResourceData) *client.Dictionary {
	maxLength := ""
	if v := data.Get(dictionaryMaxLength).(int); v > 0 {
		maxLength = strconv.Itoa(v)
	}
	choice := strconv.Itoa(data.Get(dictionaryChoice).(int))
	dictionary := client.Dictionary{
		Name:             data.Get(dictionaryName).(string),
		Element:          data.Get(dictionaryElement).(string),
		ColumnLabel:      data.Get(dictionaryColumnLabel).(string),
		InternalType:     data.Get(dictionaryInternalType).(string),
		MaxLength:        maxLength,
		Mandatory:        data.Get(dictionaryMandatory).(bool),
		ReadOnly:         data.Get(dictionaryReadOnly).(bool),
		Active:           data.Get(dictionaryActive).(bool),
		Display:          data.Get(dictionaryDisplay).(bool),
		Unique:           data.Get(dictionaryUnique).(bool),
		DefaultValue:     data.Get(dictionaryDefaultValue).(string),
		Comments:         data.Get(dictionaryComments).(string),
		Reference:        data.Get(dictionaryReference).(string),
		DynamicCreation:  data.Get(dictionaryDynamicCreation).(bool),
		Dependent:        data.Get(dictionaryDependent).(string),
		DependentOnField: data.Get(dictionaryDependentOnField).(string),
		Choice:           choice,
	}
	dictionary.ID = data.Id()
	dictionary.Scope = data.Get(commonScope).(string)
	return &dictionary
}
