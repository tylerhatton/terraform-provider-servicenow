package client

// EndpointUser is the endpoint to manage user records.
const EndpointUser = "sys_user.do"

// User is the json response for a user in ServiceNow.
type User struct {
	BaseResult
	UserName           string `json:"user_name"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email"`
	Title              string `json:"title"`
	Phone              string `json:"phone"`
	MobilePhone        string `json:"mobile_phone"`
	Active             bool   `json:"active,string"`
	LockedOut          bool   `json:"locked_out,string"`
	PasswordNeedsReset bool   `json:"password_needs_reset,string"`
	VIP                bool   `json:"vip,string"`
	TimeZone           string `json:"time_zone"`
	Location           string `json:"location"`
	Department         string `json:"department"`
	Manager            string `json:"manager"`
	Company            string `json:"company"`
	UserPassword       string `json:"user_password,omitempty"`
}
