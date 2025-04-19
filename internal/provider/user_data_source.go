package provider

import (
	"context"
	"fmt"

	"github.com/oliver-binns/googleplay-go"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

type userDataSourceModel struct {
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

type userDataSource struct {
	client *googleplay.Client
}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (d *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches user data from the Google Play API.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the user",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address for the user",
				Required:            true,
			},
		},
	}
}

func (d *userDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*googleplay.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *googleplay.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	email := data.Email.ValueString()
	if email == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'email' is required to fetch user information.",
		)
		return
	}

	// List users from Google Play API
	users, err := d.client.ListUsers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to fetch users",
			err.Error(),
		)
		return
	}

	for _, user := range users {
		if user.Email == email {
			data.Name = types.StringValue(user.Name)
			data.Email = types.StringValue(user.Email)
			break
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
