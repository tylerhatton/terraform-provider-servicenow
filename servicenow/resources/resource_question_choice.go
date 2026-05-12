package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const questionChoiceText = "text"
const questionChoiceValue = "value"
const questionChoiceQuestion = "question"
const questionChoiceOrder = "order"
const questionChoicePrice = "price"
const questionChoiceRecurringPrice = "recurring_price"
const questionChoiceInactive = "inactive"

// ResourceQuestionChoice manages a Question Choice in ServiceNow.
func ResourceQuestionChoice() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_question_choice` manages a question choice within ServiceNow.",

		CreateContext: createResourceQuestionChoice,
		ReadContext:   readResourceQuestionChoice,
		UpdateContext: updateResourceQuestionChoice,
		DeleteContext: deleteResourceQuestionChoice,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			questionChoiceText: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display text of the question choice.",
			},
			questionChoiceValue: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value of the question choice.",
			},
			questionChoiceQuestion: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The sys id of the variable/question the question choice is assigned to.",
			},
			questionChoiceOrder: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "100",
				Description: "The order the question choice will be displayed in a list of question choices.",
			},
			questionChoicePrice: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The price associated with question choice.",
			},
			questionChoiceRecurringPrice: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The recurring price associated with question choice.",
			},
			questionChoiceInactive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     "false",
				Description: "The recurring price associated with question choice.",
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceQuestionChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	questionChoice := &client.QuestionChoice{}
	if err := snowClient.GetObject(ctx, client.EndpointQuestionChoice, data.Id(), questionChoice); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromQuestionChoice(data, questionChoice)

	return nil
}

func createResourceQuestionChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	questionChoice := resourceToQuestionChoice(data)
	if err := snowClient.CreateObject(ctx, client.EndpointQuestionChoice, questionChoice); err != nil {
		return diag.FromErr(err)
	}

	resourceFromQuestionChoice(data, questionChoice)

	return readResourceQuestionChoice(ctx, data, serviceNowClient)
}

func updateResourceQuestionChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointQuestionChoice, resourceToQuestionChoice(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceQuestionChoice(ctx, data, serviceNowClient)
}

func deleteResourceQuestionChoice(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointQuestionChoice, data.Id()))
}

func resourceFromQuestionChoice(data *schema.ResourceData, questionChoice *client.QuestionChoice) {
	data.SetId(questionChoice.ID)
	data.Set(questionChoiceText, questionChoice.Text)
	data.Set(questionChoiceValue, questionChoice.Value)
	data.Set(questionChoiceQuestion, questionChoice.Question)
	data.Set(questionChoiceOrder, questionChoice.Order)
	data.Set(questionChoicePrice, questionChoice.Price)
	data.Set(questionChoiceRecurringPrice, questionChoice.RecurringPrice)
	data.Set(questionChoiceInactive, questionChoice.Inactive)
	data.Set(commonScope, questionChoice.Scope)
}

func resourceToQuestionChoice(data *schema.ResourceData) *client.QuestionChoice {
	questionChoice := client.QuestionChoice{
		Text:           data.Get(questionChoiceText).(string),
		Value:          data.Get(questionChoiceValue).(string),
		Question:       data.Get(questionChoiceQuestion).(string),
		Order:          data.Get(questionChoiceOrder).(string),
		Price:          data.Get(questionChoicePrice).(string),
		RecurringPrice: data.Get(questionChoiceRecurringPrice).(string),
		Inactive:       data.Get(questionChoiceInactive).(bool),
	}
	questionChoice.ID = data.Id()
	questionChoice.Scope = data.Get(commonScope).(string)
	return &questionChoice
}
