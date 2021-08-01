package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

		Create: createResourceQuestionChoice,
		Read:   readResourceQuestionChoice,
		Update: updateResourceQuestionChoice,
		Delete: deleteResourceQuestionChoice,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Default:     "0",
				Description: "The price associated with question choice.",
			},
			questionChoiceRecurringPrice: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
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

func readResourceQuestionChoice(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	questionChoice := &client.QuestionChoice{}
	if err := snowClient.GetObject(client.EndpointQuestionChoice, data.Id(), questionChoice); err != nil {
		data.SetId("")
		return err
	}

	resourceFromQuestionChoice(data, questionChoice)

	return nil
}

func createResourceQuestionChoice(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	questionChoice := resourceToQuestionChoice(data)
	if err := snowClient.CreateObject(client.EndpointQuestionChoice, questionChoice); err != nil {
		return err
	}

	resourceFromQuestionChoice(data, questionChoice)

	return readResourceQuestionChoice(data, serviceNowClient)
}

func updateResourceQuestionChoice(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(client.EndpointQuestionChoice, resourceToQuestionChoice(data)); err != nil {
		return err
	}

	return readResourceQuestionChoice(data, serviceNowClient)
}

func deleteResourceQuestionChoice(data *schema.ResourceData, serviceNowClient interface{}) error {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return snowClient.DeleteObject(client.EndpointQuestionChoice, data.Id())
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
