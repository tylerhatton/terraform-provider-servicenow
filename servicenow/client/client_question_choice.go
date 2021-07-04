package client

// EndpointQuestionChoice is the endpoint to manage question choice for service catalog
const EndpointQuestionChoice = "question_choice.do"

// QuestionChoice is the json response for a question choice in ServiceNow.
type QuestionChoice struct {
	BaseResult
	Text           string `json:"text"`
	Value          string `json:"value"`
	Question       string `json:"question"`
	Order          string `json:"order"`
	Price          string `json:"misc"`
	RecurringPrice string `json:"rec_misc"`
	Inactive       bool   `json:"inactive,string"`
}
