package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow/client"
)

const oauthEntityName = "name"
const oauthEntityClientUUID = "client_uuid"
const oauthEntityClientID = "client_id"
const oauthEntityAccessTokenLifespan = "access_token_lifespan"
const oauthEntityRefreshTokenLifespan = "refresh_token_lifespan"
const oauthEntityRedirectURL = "redirect_url"
const oauthEntityLogoURL = "logo_url"
const oauthEntityAccess = "access"

// ResourceOAuthEntity manages an OAuthEntity in ServiceNow.
func ResourceOAuthEntity() *schema.Resource {
	return &schema.Resource{
		Description: "`servicenow_oauth_entity` manages an OAuth application entity within ServiceNow.",

		CreateContext: createResourceOAuthEntity,
		ReadContext:   readResourceOAuthEntity,
		UpdateContext: updateResourceOAuthEntity,
		DeleteContext: deleteResourceOAuthEntity,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			oauthEntityName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the OAuth app.",
			},
			oauthEntityClientUUID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal unique identifier of the entity.",
			},
			oauthEntityClientID: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "OAuth Client ID required during handshake.",
			},
			oauthEntityAccessTokenLifespan: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Number of seconds a newly created access token is good for.",
			},
			oauthEntityRefreshTokenLifespan: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     8640000,
				Description: "Number of seconds the refresh token is good for.",
			},
			oauthEntityRedirectURL: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.Any(validation.StringIsEmpty, validation.IsURLWithScheme([]string{"http", "https"})),
				Description:  "The OAuth app's endpoint to receive the authorization code.",
			},
			oauthEntityLogoURL: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.Any(validation.StringIsEmpty, validation.IsURLWithScheme([]string{"http", "https"})),
				Description:  "URL of the logo image to display for the OAuth application.",
			},
			oauthEntityAccess: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "public",
				Description:  "Whether this Script can be accessed from only this application scope or all application scopes. Values can be 'package_private' or 'public'.",
				ValidateFunc: validation.StringInSlice([]string{"package_private", "public"}, false),
			},
			commonScope: getScopeSchema(),
		},
	}
}

func readResourceOAuthEntity(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	oauthEntity := &client.OAuthEntity{}
	if err := snowClient.GetObject(ctx, client.EndpointOAuthEntity, data.Id(), oauthEntity); err != nil {
		if client.IsNotFound(err) {
			data.SetId("")
			return nil
		}
		data.SetId("")
		return diag.FromErr(err)
	}

	resourceFromOAuthEntity(data, oauthEntity)

	return nil
}

func createResourceOAuthEntity(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	oauthEntity := resourceToOAuthEntity(data)
	if err := snowClient.CreateObject(ctx, client.EndpointOAuthEntity, oauthEntity); err != nil {
		return diag.FromErr(err)
	}

	resourceFromOAuthEntity(data, oauthEntity)

	return readResourceOAuthEntity(ctx, data, serviceNowClient)
}

func updateResourceOAuthEntity(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	if err := snowClient.UpdateObject(ctx, client.EndpointOAuthEntity, resourceToOAuthEntity(data)); err != nil {
		return diag.FromErr(err)
	}

	return readResourceOAuthEntity(ctx, data, serviceNowClient)
}

func deleteResourceOAuthEntity(ctx context.Context, data *schema.ResourceData, serviceNowClient interface{}) diag.Diagnostics {
	snowClient := serviceNowClient.(client.ServiceNowClient)
	return diag.FromErr(snowClient.DeleteObject(ctx, client.EndpointOAuthEntity, data.Id()))
}

func resourceFromOAuthEntity(data *schema.ResourceData, oauthEntity *client.OAuthEntity) {
	data.SetId(oauthEntity.ID)
	data.Set(oauthEntityName, oauthEntity.Name)
	data.Set(oauthEntityClientUUID, oauthEntity.ClientUUID)
	data.Set(oauthEntityClientID, oauthEntity.ClientID)
	data.Set(oauthEntityAccessTokenLifespan, oauthEntity.AccessTokenLifespan)
	data.Set(oauthEntityRefreshTokenLifespan, oauthEntity.RefreshTokenLifespan)
	data.Set(oauthEntityRedirectURL, oauthEntity.RedirectURL)
	data.Set(oauthEntityLogoURL, oauthEntity.LogoURL)
	data.Set(oauthEntityAccess, oauthEntity.Access)
	data.Set(commonScope, oauthEntity.Scope)
}

func resourceToOAuthEntity(data *schema.ResourceData) *client.OAuthEntity {
	oauthEntity := client.OAuthEntity{
		Name:                 data.Get(oauthEntityName).(string),
		AccessTokenLifespan:  data.Get(oauthEntityAccessTokenLifespan).(int),
		RefreshTokenLifespan: data.Get(oauthEntityRefreshTokenLifespan).(int),
		RedirectURL:          data.Get(oauthEntityRedirectURL).(string),
		LogoURL:              data.Get(oauthEntityLogoURL).(string),
		Access:               data.Get(oauthEntityAccess).(string),
	}
	oauthEntity.ID = data.Id()
	oauthEntity.Scope = data.Get(commonScope).(string)
	return &oauthEntity
}
