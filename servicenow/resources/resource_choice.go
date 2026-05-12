package resources

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const choiceName = "name"
const choiceElement = "element"
const choiceValue = "value"
const choiceLabel = "label"
const choiceSequence = "sequence"
const choiceHint = "hint"
const choiceInactive = "inactive"
const choiceDependentValue = "dependent_value"
const choiceLanguage = "language"

// ResourceChoice manages a choice list value (sys_choice entry) in ServiceNow.
func ResourceChoice() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_choice` manages a single choice list value (sys_choice entry) for a column on a table in ServiceNow.",

		CreateContext: createResourceChoice,
		ReadContext:   readResourceChoice,
		UpdateContext: updateResourceChoice,
		DeleteContext: deleteResourceChoice,

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
			choiceName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of the table this choice list value belongs to.",
			},
			choiceElement: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of the column (field) this choice list value belongs to.",
			},
			choiceValue: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The stored value for this choice.",
			},
			choiceLabel: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The human readable label displayed for this choice.",
			},
			choiceSequence: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Display order for this choice. Lower numbers appear first.",
			},
			choiceHint: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Tool tip text displayed when hovering over this choice.",
			},
			choiceInactive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, this choice is inactive and will not be displayed to users.",
			},
			choiceDependentValue: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "When this choice belongs to a dependent field, the value of the parent choice this depends on.",
			},
			choiceLanguage: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "en",
				Description: "The language code this choice's label is written in.",
			},
		},
	}
}

func readResourceChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	choice := &client.Choice{}
	if err := snowClient.GetObject(ctx, client.EndpointChoice, data.Id(), choice); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromChoice(data, choice)

	return nil
}

func createResourceChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	choice := resourceToChoice(data)
	if err := snowClient.CreateObject(ctx, client.EndpointChoice, choice); err != nil {
		return diag.FromErr(err)
	}

	resourceFromChoice(data, choice)

	return readResourceChoice(ctx, data, serviceNowClient)
}

func updateResourceChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointChoice, resourceToChoice(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceChoice(ctx, data, serviceNowClient)
}

func deleteResourceChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointChoice, data.Id()))
}

func resourceFromChoice(data *schema.ResourceData, choice *client.Choice) {
	data.SetId(choice.ID)
	data.Set(choiceName, choice.Name)
	data.Set(choiceElement, choice.Element)
	data.Set(choiceValue, choice.Value)
	data.Set(choiceLabel, choice.Label)
	sequence, err := strconv.Atoi(choice.Sequence)
	if err != nil {
		sequence = 100
	}
	data.Set(choiceSequence, sequence)
	data.Set(choiceHint, choice.Hint)
	data.Set(choiceInactive, choice.Inactive)
	data.Set(choiceDependentValue, choice.DependentValue)
	data.Set(choiceLanguage, choice.Language)
}

func resourceToChoice(data *schema.ResourceData) *client.Choice {
	choice := client.Choice{
		Name:           data.Get(choiceName).(string),
		Element:        data.Get(choiceElement).(string),
		Value:          data.Get(choiceValue).(string),
		Label:          data.Get(choiceLabel).(string),
		Sequence:       strconv.Itoa(data.Get(choiceSequence).(int)),
		Hint:           data.Get(choiceHint).(string),
		Inactive:       data.Get(choiceInactive).(bool),
		DependentValue: data.Get(choiceDependentValue).(string),
		Language:       data.Get(choiceLanguage).(string),
	}
	choice.ID = data.Id()
	return &choice
}
