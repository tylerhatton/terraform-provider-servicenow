package client

// EndpointBasicAuthCredential is the endpoint to manage basic auth credential records.
const EndpointBasicAuthCredential = "basic_auth_credentials.do"

// BasicAuthCredential is the json response for a basic auth credential in ServiceNow.
type BasicAuthCredential struct {
	BaseResult
	Name            string `json:"name"`
	Order           string `json:"order"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	CredentialAlias string `json:"tag"`
}
