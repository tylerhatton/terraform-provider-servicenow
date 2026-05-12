package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const userUserName = "user_name"
const userFirstName = "first_name"
const userLastName = "last_name"
const userEmail = "email"
const userTitle = "title"
const userPhone = "phone"
const userMobilePhone = "mobile_phone"
const userActive = "active"
const userLockedOut = "locked_out"
const userPasswordNeedsReset = "password_needs_reset"
const userVIP = "vip"
const userTimeZone = "time_zone"
const userLocation = "location"
const userDepartment = "department"
const userManager = "manager"
const userCompany = "company"
const userUserPassword = "user_password"

// ResourceUser manages a User in ServiceNow.
func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_user` manages a user record within ServiceNow.",

		CreateContext: createResourceUser,
		ReadContext:   readResourceUser,
		UpdateContext: updateResourceUser,
		DeleteContext: deleteResourceUser,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			userUserName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The login user name for the user.",
			},
			userFirstName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "First name of the user.",
			},
			userLastName: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Last name of the user.",
			},
			userEmail: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Email address of the user.",
			},
			userTitle: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Business title of the user.",
			},
			userPhone: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Business phone number of the user.",
			},
			userMobilePhone: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Mobile phone number of the user.",
			},
			userActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not this user is active and able to log in.",
			},
			userLockedOut: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the user account is locked out.",
			},
			userPasswordNeedsReset: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the user must reset their password at next login.",
			},
			userVIP: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the user is flagged as a VIP.",
			},
			userTimeZone: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Time zone preference for the user (for example 'US/Eastern').",
			},
			userLocation: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the location record the user is associated with.",
			},
			userDepartment: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the department record the user belongs to.",
			},
			userManager: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the user record that is this user's manager.",
			},
			userCompany: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Sys ID of the company record the user belongs to.",
			},
			userUserPassword: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "Plain text password for the user. ServiceNow stores this value hashed and may not return the original; this field is write-only and not refreshed from state.",
			},
		},
	}
}

func readResourceUser(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	user := &client.User{}
	if err := snowClient.GetObject(ctx, client.EndpointUser, data.Id(), user); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromUser(data, user)

	return nil
}

func createResourceUser(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	user := resourceToUser(data)
	if err := snowClient.CreateObject(ctx, client.EndpointUser, user); err != nil {
		return diag.FromErr(err)
	}

	resourceFromUser(data, user)

	return readResourceUser(ctx, data, serviceNowClient)
}

func updateResourceUser(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointUser, resourceToUser(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceUser(ctx, data, serviceNowClient)
}

func deleteResourceUser(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointUser, data.Id()))
}

func resourceFromUser(data *schema.ResourceData, user *client.User) {
	data.SetId(user.ID)
	data.Set(userUserName, user.UserName)
	data.Set(userFirstName, user.FirstName)
	data.Set(userLastName, user.LastName)
	data.Set(userEmail, user.Email)
	data.Set(userTitle, user.Title)
	data.Set(userPhone, user.Phone)
	data.Set(userMobilePhone, user.MobilePhone)
	data.Set(userActive, user.Active)
	data.Set(userLockedOut, user.LockedOut)
	data.Set(userPasswordNeedsReset, user.PasswordNeedsReset)
	data.Set(userVIP, user.VIP)
	data.Set(userTimeZone, user.TimeZone)
	data.Set(userLocation, user.Location)
	data.Set(userDepartment, user.Department)
	data.Set(userManager, user.Manager)
	data.Set(userCompany, user.Company)
	// ServiceNow stores the password hashed and the returned value will not match the user
	// supplied plain text. Avoid overwriting the user-provided value in state.
}

func resourceToUser(data *schema.ResourceData) *client.User {
	user := client.User{
		UserName:           data.Get(userUserName).(string),
		FirstName:          data.Get(userFirstName).(string),
		LastName:           data.Get(userLastName).(string),
		Email:              data.Get(userEmail).(string),
		Title:              data.Get(userTitle).(string),
		Phone:              data.Get(userPhone).(string),
		MobilePhone:        data.Get(userMobilePhone).(string),
		Active:             data.Get(userActive).(bool),
		LockedOut:          data.Get(userLockedOut).(bool),
		PasswordNeedsReset: data.Get(userPasswordNeedsReset).(bool),
		VIP:                data.Get(userVIP).(bool),
		TimeZone:           data.Get(userTimeZone).(string),
		Location:           data.Get(userLocation).(string),
		Department:         data.Get(userDepartment).(string),
		Manager:            data.Get(userManager).(string),
		Company:            data.Get(userCompany).(string),
		UserPassword:       data.Get(userUserPassword).(string),
	}
	user.ID = data.Id()
	return &user
}
