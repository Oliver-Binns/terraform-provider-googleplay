package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	androidpublisher "google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

var _ provider.Provider = &GooglePlayProvider{}

// GooglePlayClient wraps the official Android Publisher service with a developer ID.
type GooglePlayClient struct {
	service     *androidpublisher.Service
	developerID string
}

func (c *GooglePlayClient) ListUsers(ctx context.Context) ([]*androidpublisher.User, error) {
	parent := fmt.Sprintf("developers/%s", c.developerID)
	resp, err := c.service.Users.List(parent).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func (c *GooglePlayClient) CreateUser(
	email string,
	permissions []DeveloperLevelPermission,
	ctx context.Context,
) (*androidpublisher.User, error) {
	parent := fmt.Sprintf("developers/%s", c.developerID)
	perms := make([]string, len(permissions))
	for i, p := range permissions {
		perms[i] = string(p)
	}
	user := &androidpublisher.User{
		Email:                       email,
		DeveloperAccountPermissions: perms,
	}
	return c.service.Users.Create(parent, user).Context(ctx).Do()
}

func (c *GooglePlayClient) UpdateUser(
	email string,
	permissions *[]DeveloperLevelPermission,
	ctx context.Context,
) (*androidpublisher.User, error) {
	name := fmt.Sprintf("developers/%s/users/%s", c.developerID, email)
	perms := make([]string, len(*permissions))
	for i, p := range *permissions {
		perms[i] = string(p)
	}
	user := &androidpublisher.User{
		DeveloperAccountPermissions: perms,
	}
	return c.service.Users.Patch(name, user).UpdateMask("developerAccountPermissions").Context(ctx).Do()
}

func (c *GooglePlayClient) DeleteUser(email string, ctx context.Context) error {
	name := fmt.Sprintf("developers/%s/users/%s", c.developerID, email)
	return c.service.Users.Delete(name).Context(ctx).Do()
}

func (c *GooglePlayClient) GrantAccess(
	email string,
	appID string,
	permissions []AppLevelPermission,
	ctx context.Context,
) (*androidpublisher.Grant, error) {
	parent := fmt.Sprintf("developers/%s/users/%s", c.developerID, email)
	perms := make([]string, len(permissions))
	for i, p := range permissions {
		perms[i] = string(p)
	}
	grant := &androidpublisher.Grant{
		PackageName:         appID,
		AppLevelPermissions: perms,
	}
	return c.service.Grants.Create(parent, grant).Context(ctx).Do()
}

func (c *GooglePlayClient) ModifyAccess(
	email string,
	appID string,
	permissions []AppLevelPermission,
	ctx context.Context,
) (*androidpublisher.Grant, error) {
	name := fmt.Sprintf("developers/%s/users/%s/grants/%s", c.developerID, email, appID)
	perms := make([]string, len(permissions))
	for i, p := range permissions {
		perms[i] = string(p)
	}
	grant := &androidpublisher.Grant{
		AppLevelPermissions: perms,
	}
	return c.service.Grants.Patch(name, grant).UpdateMask("appLevelPermissions").Context(ctx).Do()
}

func (c *GooglePlayClient) RevokeAccess(
	email string,
	appID string,
	ctx context.Context,
) error {
	name := fmt.Sprintf("developers/%s/users/%s/grants/%s", c.developerID, email, appID)
	return c.service.Grants.Delete(name).Context(ctx).Do()
}

// GooglePlayProvider defines the provider implementation.
type GooglePlayProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type GooglePlayProviderModel struct {
	ServiceAccountJson types.String `tfsdk:"service_account_json_base64"`
	DeveloperID        types.String `tfsdk:"developer_id"`
}

func (p *GooglePlayProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "googleplay"
	resp.Version = p.version
}

func (p *GooglePlayProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Google Play Console",
		Attributes: map[string]schema.Attribute{
			"service_account_json_base64": schema.StringAttribute{
				MarkdownDescription: `The service account JSON data used to authenticate with Google:
				https://developers.google.com/android-publisher/getting_started#service-account`,
				Optional:  true,
				Sensitive: true,
			},
			"developer_id": schema.StringAttribute{
				MarkdownDescription: `Your unique 19-digit Google Play Developer account ID:
				https://support.google.com/googleplay/android-developer/answer/13634081?hl=en-GB`,
				Required:  true,
				Sensitive: false,
			},
		},
	}
}

func (p *GooglePlayProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data GooglePlayProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	developerID := data.DeveloperID.ValueString()

	google_credentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if !data.ServiceAccountJson.IsNull() && !data.ServiceAccountJson.IsUnknown() {
		tflog.Info(ctx, "Using service account from provider configuration")

		serviceAccountBase64 := data.ServiceAccountJson.ValueString()
		rawJson, err := base64.StdEncoding.DecodeString(serviceAccountBase64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing service account JSON",
				err.Error(),
			)
			return
		}

		service, err := androidpublisher.NewService(ctx, option.WithAuthCredentialsJSON(option.ServiceAccount, rawJson))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Android Publisher service",
				err.Error(),
			)
			return
		}

		client := &GooglePlayClient{service: service, developerID: developerID}
		resp.DataSourceData = client
		resp.ResourceData = client
	} else if google_credentials != "" {
		tflog.Info(ctx, "Using service account from GOOGLE_APPLICATION_CREDENTIALS environment variable")

		service, err := androidpublisher.NewService(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Android Publisher service",
				err.Error(),
			)
			return
		}

		tflog.Info(ctx, "created client successfully")

		client := &GooglePlayClient{service: service, developerID: developerID}
		resp.DataSourceData = client
		resp.ResourceData = client

	} else {
		resp.Diagnostics.AddError(
			"Missing service account JSON",
			"Please either provide service account credentials either in the provider configuration or set the GOOGLE_APPLICATION_CREDENTIALS environment variable.",
		)
		return
	}
}

func (p *GooglePlayProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewAppIAMResource,
	}
}

func (p *GooglePlayProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *GooglePlayProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GooglePlayProvider{
			version: version,
		}
	}
}
