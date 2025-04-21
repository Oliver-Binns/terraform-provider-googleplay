package provider

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/oliver-binns/googleplay-go"
)

var _ provider.Provider = &GooglePlayProvider{}

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
				Required:  true,
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

	serviceAccountBase64 := data.ServiceAccountJson.ValueString()
	rawJson, err := base64.StdEncoding.DecodeString(serviceAccountBase64)
	log(err, resp)

	developerID := data.DeveloperID.ValueString()

	client := googleplay.GooglePlayClient(
		developerID,
		string(rawJson),
	)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func log(err error, resp *provider.ConfigureResponse) {
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing service account JSON",
			err.Error(),
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
